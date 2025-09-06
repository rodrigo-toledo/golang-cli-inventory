package testutils

import (
	"cli-inventory/internal/db"
	"context"
	"fmt"
	"log"
	"os"
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
func SetupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	once.Do(func() {
		var err error

		// Create a pool of Docker clients
		testPool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}

		// Pull the PostgreSQL image
		err = testPool.Client.PullImage(docker.PullImageOptions{
			Repository: "postgres",
			Tag:        "15",
		}, docker.AuthConfiguration{})
		if err != nil {
			log.Fatalf("Could not pull PostgreSQL image: %s", err)
		}

		// Create a container with PostgreSQL
		testResource, err = testPool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15",
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
		databaseURL := fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", hostAndPort)

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

	return testDB
}

// TeardownTestDatabase stops and removes the test database container
func TeardownTestDatabase(t *testing.T) {
	t.Helper()

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

	// Disable foreign key checks temporarily
	_, err := db.Exec(ctx, "SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		t.Fatalf("Could not disable constraints: %s", err)
	}

	// Truncate all tables in the correct order to respect foreign key constraints
	tables := []string{
		"stock_movements",
		"stock",
		"products",
		"locations",
	}

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Fatalf("Could not truncate table %s: %s", table, err)
		}
	}

	// Reset sequences
	sequences := []string{
		"locations_id_seq",
		"products_id_seq",
		"stock_id_seq",
		"stock_movements_id_seq",
	}

	for _, seq := range sequences {
		_, err := db.Exec(ctx, fmt.Sprintf("ALTER SEQUENCE %s RESTART WITH 1", seq))
		if err != nil {
			t.Fatalf("Could not reset sequence %s: %s", seq, err)
		}
	}
}

// runMigrations creates the database schema for testing
func runMigrations(db *pgxpool.Pool) error {
	ctx := context.Background()

	// Create tables schema
	schema := `
	CREATE TABLE locations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE products (
		id SERIAL PRIMARY KEY,
		sku VARCHAR(50) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10, 2),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE stock (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		location_id INTEGER NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
		quantity INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(product_id, location_id)
	);

	CREATE TABLE stock_movements (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		from_location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
		to_location_id INTEGER REFERENCES locations(id) ON DELETE SET NULL,
		quantity INTEGER NOT NULL,
		movement_type VARCHAR(50) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`

	_, err := db.Exec(ctx, schema)
	if err != nil {
		return fmt.Errorf("could not create tables: %w", err)
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
