package database

import (
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