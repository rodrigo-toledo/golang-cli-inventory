package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestCreateJWT(t *testing.T) {
	user := &User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}
	
	secret := "test-secret"
	expirationTime := time.Now().Add(1 * time.Hour)
	
	tokenString, err := CreateJWT(user, secret, expirationTime)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	
	// Parse and verify the token
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	
	assert.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, user.Name, claims.Name)
}

func TestAuthenticator_NoAuthHeader(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)
	
	// Create a test request without authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	
	// Create a simple handler that should not be called
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})
	
	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)
	
	// Check response
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Authorization header or session token cookie required")
}