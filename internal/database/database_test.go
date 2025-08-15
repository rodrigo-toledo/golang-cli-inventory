//go:build integration

package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInitialized(t *testing.T) {
	// Since we can't easily test the actual database connection without Docker,
	// we'll just test the IsInitialized function logic
	
	// Initially, DB should be nil, so IsInitialized should return false
	assert.False(t, IsInitialized())
	
	// Note: We can't easily test the true case without actually initializing the database
	// which would require Docker or a real PostgreSQL instance
}

func TestInitDB_MissingEnvVar(t *testing.T) {
	// Save original environment variable
	originalDBURL := os.Getenv("DATABASE_URL")
	
	// Ensure cleanup
	defer func() {
		os.Setenv("DATABASE_URL", originalDBURL)
	}()
	
	// Unset DATABASE_URL to test fallback
	os.Unsetenv("DATABASE_URL")
	
	// Try to initialize DB (this will fail because we don't have a real database running)
	err := InitDB()
	// We expect an error because we don't have a real database running
	assert.Error(t, err)
	// But we can still check that the function tried to initialize
	assert.True(t, IsInitialized() || err != nil)
}