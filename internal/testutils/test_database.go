package testutils

import (
	"cli-inventory/internal/db"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

// SetupTestDatabase creates a test database using Docker and returns the connection pool.
// It ensures a clean database for each test that calls it.
func SetupTestDatabase(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// Use dockertest to create a new postgres container for each test.
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	// Define the network name from docker-compose
	networkName := "golang-cli-inventory_default"

	// Check if the network exists
	networks, err := pool.Client.ListNetworks()
	if err != nil {
		t.Fatalf("Could not list Docker networks: %s", err)
	}

	var network *docker.Network
	for i := range networks {
		if networks[i].Name == networkName {
			network = &networks[i]
			break
		}
	}

	if network == nil {
		t.Fatalf("Docker network '%s' not found. Please ensure docker-compose is running.", networkName)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpass",
			"POSTGRES_DB=testdb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	// Connect the container to the network
	err = pool.Client.ConnectNetwork(network.ID, docker.NetworkConnectionOptions{
		Container: resource.Container.ID,
	})
	if err != nil {
		t.Fatalf("Could not connect container to network '%s': %s", networkName, err)
	}

	// Set a timeout for the container to be purged.
	resource.Expire(120) // 2 minutes

	// Create a database connection pool.
	// Use the container name as the hostname, which is resolvable within the docker network.
	var dbPool *pgxpool.Pool
	containerName := resource.Container.Name[1:] // Remove the leading '/'
	databaseURL := fmt.Sprintf("postgres://testuser:testpass@%s:5432/testdb?sslmode=disable", containerName)

	// Wait for the database to be ready.
	if err := pool.Retry(func() error {
		var err error
		dbPool, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return dbPool.Ping(context.Background())
	}); err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	// Run migrations.
	if err := runMigrations(dbPool); err != nil {
		t.Fatalf("Could not run migrations: %s", err)
	}

	// Return the connection pool and a cleanup function.
	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Errorf("Could not purge resource: %s", err)
		}
	})

	return dbPool
}

// TeardownTestDatabase is no longer needed as t.Cleanup is used in SetupTestDatabase.
//
// Deprecated: use t.Cleanup via SetupTestDatabase instead.
func TeardownTestDatabase(_ *testing.T) {
}

// CleanupTestDatabase truncates all tables between tests
func CleanupTestDatabase(t *testing.T, pool *pgxpool.Pool) {
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
		_, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table))
		if err != nil {
			t.Fatalf("Could not truncate table %s: %s", table, err)
		}
	}
}

// runMigrations creates the database schema for testing
func runMigrations(pool *pgxpool.Pool) error {
	// Read the migration file
	migration, err := os.ReadFile("../../migrations/000001_create_tables.up.sql")
	if err != nil {
		return fmt.Errorf("could not read migration file: %w", err)
	}

	// Execute the migration
	_, err = pool.Exec(context.Background(), string(migration))
	if err != nil {
		// Do not ignore "already exists" errors, as this indicates a problem with test isolation.
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

// GetTestQueries creates and returns a new Queries instance for testing
func GetTestQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
