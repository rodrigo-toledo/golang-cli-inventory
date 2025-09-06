package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserFromContext(t *testing.T) {
	// Test that UserFromContext returns false when no user is in context
	user, ok := UserFromContext(context.Background())
	assert.Nil(t, user)
	assert.False(t, ok)
}

func TestSessionSecret(t *testing.T) {
	// Create a minimal AuthHandler to test SessionSecret method
	handler := &AuthHandler{
		sessionSecret: "test-secret",
	}
	
	assert.Equal(t, "test-secret", handler.SessionSecret())
}