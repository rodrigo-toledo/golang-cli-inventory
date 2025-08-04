// Package database provides database connection functionality for the inventory management system.
// It handles the initialization and management of the PostgreSQL database connection pool.
package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool that can be used throughout the application.
// It is initialized by calling InitDB().
var DB *pgxpool.Pool

var once sync.Once
var initErr error

// InitDB initializes the database connection pool using a singleton pattern.
// It reads the database connection URL from the DATABASE_URL environment variable,
// falling back to a default local development URL if the environment variable is not set.
// It also tests the connection to ensure it's working properly.
// Returns an error if the connection fails, allowing for graceful handling.
func InitDB() error {
	var err error
	once.Do(func() {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			// Fallback to default for local development
			databaseURL = "postgres://inventory_user:inventory_password@localhost:5432/inventory_db?sslmode=disable"
		}

		DB, err = pgxpool.New(context.Background(), databaseURL)
		if err != nil {
			initErr = fmt.Errorf("unable to connect to database: %w", err)
			return
		}

		// Test the connection
		err = DB.Ping(context.Background())
		if err != nil {
			initErr = fmt.Errorf("unable to ping database: %w", err)
			return
		}

		fmt.Println("Connected to database successfully")
	})

	return initErr
}

// IsInitialized checks if the database connection has been initialized
func IsInitialized() bool {
	return DB != nil
}
