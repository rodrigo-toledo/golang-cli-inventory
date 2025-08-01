package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

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
