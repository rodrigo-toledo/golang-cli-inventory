package testutils

import (
	"cli-inventory/internal/db"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	testDB       *pgxpool.Pool
	testPool     *dockertest.Pool
	testResource *dockertest.Resource
	once         sync.Once
)

// SetupTestDatabase creates a test database using Docker and returns the connection pool
// This function uses dockertest to manage the container lifecycle
// If DATABASE_URL is set, it will use that connection instead of creating a new container
func SetupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Check if we're running in a Docker environment with an existing database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		once.Do(func() {
			// Connect to the default 'postgres' database to create the test database
			var defaultDb *pgxpool.Pool
			var err error
			deadline := time.Now().Add(60 * time.Second)
			for {
				defaultDb, err = pgxpool.New(context.Background(), "postgres://inventory_user:inventory_password@db:5432/postgres?sslmode=disable")
				if err == nil {
					if pingErr := defaultDb.Ping(context.Background()); pingErr == nil {
						break
					} else {
						err = pingErr
						if defaultDb != nil {
							defaultDb.Close()
						}
					}
				}
				if time.Now().After(deadline) {
					log.Fatalf("Could not connect to default database after retries: %s", err)
				}
				time.Sleep(1 * time.Second)
			}

			// Create the test database if it doesn't exist
			_, err = defaultDb.Exec(context.Background(), "CREATE DATABASE inventory_test_db")
			if err != nil {
				// Ignore "already exists" error
				if !strings.Contains(err.Error(), "already exists") {
					log.Fatalf("Could not create test database: %s", err)
				}
			}
			defaultDb.Close()

			// Now, connect to the test database
			deadline = time.Now().Add(60 * time.Second)
			for {
				testDB, err = pgxpool.New(context.Background(), databaseURL)
				if err == nil {
					if pingErr := testDB.Ping(context.Background()); pingErr == nil {
						break
					} else {
						err = pingErr
						if testDB != nil {
							testDB.Close()
						}
					}
				}
				if time.Now().After(deadline) {
					log.Fatalf("Could not connect to existing database after retries: %s", err)
				}
				time.Sleep(1 * time.Second)
			}

			// Run migrations
			if err := runMigrations(testDB); err != nil {
				log.Fatalf("Could not run migrations: %s", err)
			}
		})

		// Cleanup the database immediately after connecting
		CleanupTestDatabase(t, testDB)

		return testDB
	}

	// For Docker testing, we need to ensure each test gets a clean database
	// We'll use a singleton pattern but ensure the database is clean
	once.Do(func() {
		// Create a pool of Docker clients
		var err error
		testPool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}

		// Pull the PostgreSQL image
		err = testPool.Client.PullImage(docker.PullImageOptions{
			Repository: "postgres",
			Tag:        "17",
		}, docker.AuthConfiguration{})
		if err != nil {
			log.Fatalf("Could not pull PostgreSQL image: %s", err)
		}

		// Create a container with PostgreSQL
		testResource, err = testPool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "17",
			Env: []string{
				"POSTGRES_USER=testuser",
				"POSTGRES_PASSWORD=testpass",
				"POSTGRES_DB=testdb",
				"listen_addresses = '*'",
			},
		}, func(config *docker.HostConfig) {
			// Set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}

		// Set the container to expire after 60 minutes to prevent resource leaks
		testResource.Expire(uint(60 * time.Minute / time.Second))

		// Get the database connection string
		hostAndPort := testResource.GetHostPort("5432/tcp")
		databaseURL = fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", hostAndPort)

		// Exponential backoff-retry, because the application in the container might not be ready to accept connections yet
		if err := testPool.Retry(func() error {
			var err error
			testDB, err = pgxpool.New(context.Background(), databaseURL)
			if err != nil {
				return err
			}
			return testDB.Ping(context.Background())
		}); err != nil {
			log.Fatalf("Could not connect to database: %s", err)
		}

		// Set the DATABASE_URL environment variable for the application
		os.Setenv("DATABASE_URL", databaseURL)

		// Run migrations
		if err := runMigrations(testDB); err != nil {
			log.Fatalf("Could not run migrations: %s", err)
		}
	})

	// For integration tests with Docker, we need to cleanup the database for each test
	if testDB != nil {
		CleanupTestDatabase(t, testDB)
	}

	return testDB
}

// TeardownTestDatabase stops and removes the test database container
func TeardownTestDatabase(t *testing.T) {
	t.Helper()

	// Only teardown if we're in standalone mode (DATABASE_URL not set)
	if os.Getenv("DATABASE_URL") != "" {
		return
	}

	if testResource != nil {
		if err := testPool.Purge(testResource); err != nil {
			t.Errorf("Could not purge resource: %s", err)
		}
	}
}

// CleanupTestDatabase truncates all tables between tests
func CleanupTestDatabase(t *testing.T, db *pgxpool.Pool) {
	t.Helper()

	ctx := context.Background()

	// Truncate all tables in the correct order to respect foreign key constraints
	tables := []string{
		"stock_movements",
		"stock",
		"products",
		"locations",
	}

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Fatalf("Could not truncate table %s: %s", table, err)
		}
	}
}

// runMigrations creates the database schema for testing
func runMigrations(db *pgxpool.Pool) error {
	// Read the migration file
	migration, err := os.ReadFile("../../migrations/000001_create_tables.up.sql")
	if err != nil {
		return fmt.Errorf("could not read migration file: %w", err)
	}

	// Execute the migration
	_, err = db.Exec(context.Background(), string(migration))
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("could not run migrations: %w", err)
		}
	}

	return nil
}

// GetTestDB returns the existing test database connection
func GetTestDB() *pgxpool.Pool {
	return testDB
}

// GetTestQueries creates and returns a new Queries instance for testing
func GetTestQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
