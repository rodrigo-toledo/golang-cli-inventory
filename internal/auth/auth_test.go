//go:build unit

package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestLoadConfig(t *testing.T) {
	// Save original environment variables
	originalClientID := os.Getenv("OAUTH_CLIENT_ID")
	originalClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	originalAuthURL := os.Getenv("OAUTH_AUTH_URL")
	originalTokenURL := os.Getenv("OAUTH_TOKEN_URL")
	originalRedirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	originalSessionSecret := os.Getenv("SESSION_SECRET")

	// Ensure cleanup
	defer func() {
		os.Setenv("OAUTH_CLIENT_ID", originalClientID)
		os.Setenv("OAUTH_CLIENT_SECRET", originalClientSecret)
		os.Setenv("OAUTH_AUTH_URL", originalAuthURL)
		os.Setenv("OAUTH_TOKEN_URL", originalTokenURL)
		os.Setenv("OAUTH_REDIRECT_URL", originalRedirectURL)
		os.Setenv("SESSION_SECRET", originalSessionSecret)
	}()

	// Test case 1: All required environment variables are set
	os.Setenv("OAUTH_CLIENT_ID", "test-client-id")
	os.Setenv("OAUTH_CLIENT_SECRET", "test-client-secret")
	os.Setenv("OAUTH_AUTH_URL", "https://example.com/auth")
	os.Setenv("OAUTH_TOKEN_URL", "https://example.com/token")
	os.Setenv("OAUTH_REDIRECT_URL", "https://example.com/callback")
	os.Setenv("SESSION_SECRET", "test-session-secret")

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "test-client-id", cfg.OAuthClientID)
	assert.Equal(t, "test-client-secret", cfg.OAuthClientSecret)
	assert.Equal(t, "https://example.com/auth", cfg.OAuthAuthURL)
	assert.Equal(t, "https://example.com/token", cfg.OAuthTokenURL)
	assert.Equal(t, "https://example.com/callback", cfg.OAuthRedirectURL)
	assert.Equal(t, "test-session-secret", cfg.SessionSecret)
	assert.Equal(t, []string{"openid", "profile", "email"}, cfg.OAuthScopes) // Default scopes

	// Test case 2: Custom scopes
	os.Setenv("OAUTH_SCOPES", "openid profile email custom_scope")
	cfg, err = LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"openid", "profile", "email", "custom_scope"}, cfg.OAuthScopes)

	// Test case 3: Allowed issuers
	os.Setenv("ALLOWED_ISSUERS", "https://issuer1.com, https://issuer2.com")
	cfg, err = LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"https://issuer1.com", "https://issuer2.com"}, cfg.AllowedIssuers)

	// Test case 4: Missing required environment variables
	os.Unsetenv("OAUTH_CLIENT_ID")
	cfg, err = LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestOAuth2Config(t *testing.T) {
	cfg := &AuthConfig{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		OAuthAuthURL:      "https://example.com/auth",
		OAuthTokenURL:     "https://example.com/token",
		OAuthRedirectURL:  "https://example.com/callback",
		OAuthScopes:       []string{"openid", "profile"},
	}

	oauth2Config := cfg.OAuth2Config()
	assert.Equal(t, "test-client-id", oauth2Config.ClientID)
	assert.Equal(t, "test-client-secret", oauth2Config.ClientSecret)
	assert.Equal(t, "https://example.com/callback", oauth2Config.RedirectURL)
	assert.Equal(t, []string{"openid", "profile"}, oauth2Config.Scopes)
	assert.Equal(t, "https://example.com/auth", oauth2Config.Endpoint.AuthURL)
	assert.Equal(t, "https://example.com/token", oauth2Config.Endpoint.TokenURL)
}

func TestNewAuthHandler(t *testing.T) {
	cfg := &AuthConfig{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		OAuthAuthURL:      "https://example.com/auth",
		OAuthTokenURL:     "https://example.com/token",
		OAuthRedirectURL:  "https://example.com/callback",
		SessionSecret:     "test-session-secret",
		AllowedIssuers:    []string{"https://issuer.example.com"},
	}

	handler, err := NewAuthHandler(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.Equal(t, "test-session-secret", handler.SessionSecret())
	assert.Equal(t, cfg.AllowedIssuers, handler.allowedIssuers)
}

func TestLoginHandler(t *testing.T) {
	cfg := &AuthConfig{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		OAuthAuthURL:      "https://example.com/auth",
		OAuthTokenURL:     "https://example.com/token",
		OAuthRedirectURL:  "https://example.com/callback",
		SessionSecret:     "test-session-secret",
	}

	handler, err := NewAuthHandler(cfg)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/login", nil)
	rec := httptest.NewRecorder()

	handler.LoginHandler(rec, req)

	// Check that we get a redirect
	assert.Equal(t, http.StatusFound, rec.Code)

	// Check that the state cookie is set
	cookies := rec.Result().Cookies()
	var stateCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "oauth_state" {
			stateCookie = cookie
			break
		}
	}
	assert.NotNil(t, stateCookie)
	assert.NotEmpty(t, stateCookie.Value)
	assert.Equal(t, "/", stateCookie.Path)
	assert.Equal(t, 300, stateCookie.MaxAge)
	assert.True(t, stateCookie.HttpOnly)
	assert.True(t, stateCookie.Secure)
}

func TestLogoutHandler(t *testing.T) {
	cfg := &AuthConfig{
		OAuthClientID:     "test-client-id",
		OAuthClientSecret: "test-client-secret",
		OAuthAuthURL:      "https://example.com/auth",
		OAuthTokenURL:     "https://example.com/token",
		OAuthRedirectURL:  "https://example.com/callback",
		SessionSecret:     "test-session-secret",
	}

	handler, err := NewAuthHandler(cfg)
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/logout", nil)
	rec := httptest.NewRecorder()

	handler.LogoutHandler(rec, req)

	// Check that we get a redirect
	assert.Equal(t, http.StatusFound, rec.Code)

	// Check that the session cookie is cleared
	cookies := rec.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}
	assert.NotNil(t, sessionCookie)
	assert.Empty(t, sessionCookie.Value)
	assert.Equal(t, -1, sessionCookie.MaxAge)
}