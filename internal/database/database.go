// Package database provides database connection functionality for the inventory management system.
// It handles the initialization and management of the PostgreSQL database connection pool.
package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the global database connection pool that can be used throughout the application.
// It is initialized by calling InitDB().
var DB *pgxpool.Pool

// InitDB initializes the database connection pool.
// It reads the database connection URL from the DATABASE_URL environment variable,
// falling back to a default local development URL if the environment variable is not set.
// It also tests the connection to ensure it's working properly.
func InitDB() {
	var err error
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback to default for local development
		databaseURL = "postgres://inventory_user:inventory_password@localhost:5432/inventory_db?sslmode=disable"
	}

	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Test the connection
	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	fmt.Println("Connected to database successfully")
}
