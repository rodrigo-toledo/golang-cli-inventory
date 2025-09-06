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

func TestAuthenticator_InvalidAuthHeaderFormat(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)

	// Create a test request with invalid authorization header format
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()

	// Create a simple handler that should not be called
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid authorization header format")
}

func TestAuthenticator_InvalidToken(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)

	// Create a test request with invalid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()

	// Create a simple handler that should not be called
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid token")
}

func TestAuthenticator_ValidToken(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)

	// Create a test user and token
	user := &User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	tokenString, err := CreateJWT(user, secret, expirationTime)
	assert.NoError(t, err)

	// Create a test request with valid token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	// Create a handler that checks if user is in context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userFromContext, ok := UserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, user.ID, userFromContext.ID)
		assert.Equal(t, user.Email, userFromContext.Email)
		assert.Equal(t, user.Name, userFromContext.Name)
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticator_ValidTokenFromCookie(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)

	// Create a test user and token
	user := &User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	tokenString, err := CreateJWT(user, secret, expirationTime)
	assert.NoError(t, err)

	// Create a test request with valid token in cookie
	req := httptest.NewRequest("GET", "/test", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: tokenString,
	})
	rec := httptest.NewRecorder()

	// Create a handler that checks if user is in context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userFromContext, ok := UserFromContext(r.Context())
		assert.True(t, ok)
		assert.Equal(t, user.ID, userFromContext.ID)
		assert.Equal(t, user.Email, userFromContext.Email)
		assert.Equal(t, user.Name, userFromContext.Name)
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthenticator_ExpiredToken(t *testing.T) {
	secret := "test-secret"
	middleware := Authenticator(secret)

	// Create a test user and expired token
	user := &User{
		ID:    "test-user-id",
		Email: "test@example.com",
		Name:  "Test User",
	}

	expirationTime := time.Now().Add(-1 * time.Hour) // Expired 1 hour ago
	tokenString, err := CreateJWT(user, secret, expirationTime)
	assert.NoError(t, err)

	// Create a test request with expired token
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	// Create a simple handler that should not be called
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	// Apply middleware
	middleware(handler).ServeHTTP(rec, req)

	// Check response
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Token has expired")
}