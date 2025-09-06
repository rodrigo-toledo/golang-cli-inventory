// Package auth provides authentication and authorization logic for the application.
// It includes OAuth 2.0 flow handling, session management with JWTs,
// and middleware for protecting routes.
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// AuthConfig holds the configuration for OAuth 2.0 and JWT settings.
// It is populated from environment variables.
type AuthConfig struct {
	OAuthClientID     string
	OAuthClientSecret string
	OAuthAuthURL      string
	OAuthTokenURL     string
	OAuthRedirectURL  string
	OAuthScopes       []string
	SessionSecret     string
	AllowedIssuers    []string
}

// LoadConfig loads authentication configuration from environment variables.
func LoadConfig() (*AuthConfig, error) {
	cfg := &AuthConfig{
		OAuthClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		OAuthClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		OAuthAuthURL:      os.Getenv("OAUTH_AUTH_URL"),
		OAuthTokenURL:     os.Getenv("OAUTH_TOKEN_URL"),
		OAuthRedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		SessionSecret:     os.Getenv("SESSION_SECRET"),
	}

	if cfg.OAuthClientID == "" || cfg.OAuthClientSecret == "" || cfg.OAuthAuthURL == "" ||
		cfg.OAuthTokenURL == "" || cfg.OAuthRedirectURL == "" || cfg.SessionSecret == "" {
		return nil, errors.New("missing required OAuth 2.0 environment variables")
	}

	scopes := os.Getenv("OAUTH_SCOPES")
	if scopes != "" {
		cfg.OAuthScopes = strings.Split(scopes, " ")
	} else {
		// Default scopes if not specified
		cfg.OAuthScopes = []string{"openid", "profile", "email"}
	}

	issuers := os.Getenv("ALLOWED_ISSUERS")
	if issuers != "" {
		cfg.AllowedIssuers = strings.Split(issuers, ",")
		// Trim whitespace from each issuer
		for i, issuer := range cfg.AllowedIssuers {
			cfg.AllowedIssuers[i] = strings.TrimSpace(issuer)
		}
	} else {
		// If no issuers are specified, we cannot validate tokens.
		// Depending on strictness, this could be an error or a warning.
		// For now, we'll allow it but token validation will be less strict.
		// Consider making this a required field for higher security.
	}

	return cfg, nil
}

// OAuth2Config returns a configured oauth2.Config instance.
func (c *AuthConfig) OAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.OAuthClientID,
		ClientSecret: c.OAuthClientSecret,
		RedirectURL:  c.OAuthRedirectURL,
		Scopes:       c.OAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  c.OAuthAuthURL,
			TokenURL: c.OAuthTokenURL,
		},
	}
}

// User represents the authenticated user's information.
type User struct {
	ID    string
	Email string
	Name  string
	// Add other fields as needed from the ID token
}

// contextKey is a private type for context keys to avoid collisions.
type contextKey string

const userContextKey = contextKey("user")

// UserFromContext retrieves the User object from the request context.
func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userContextKey).(*User)
	return user, ok
}

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	oauth2Config   *oauth2.Config
	provider       *oidc.Provider
	verifier       *oidc.IDTokenVerifier
	allowedIssuers []string
	sessionSecret  string
}

// NewAuthHandler creates a new AuthHandler.
// It initializes the OIDC provider and verifier.
func NewAuthHandler(cfg *AuthConfig) (*AuthHandler, error) {
	oauth2Config := cfg.OAuth2Config()

	// We need at least one issuer to set up the OIDC provider.
	// The provider's issuer will be used to fetch discovery document.
	// For generic OAuth2, OIDC might not be fully supported by all providers.
	// We'll try to use the first allowed issuer as the basis for OIDC.
	// If no issuers are configured, we fall back to a more basic OAuth2 flow
	// without ID token verification, which is less secure.

	var provider *oidc.Provider
	var verifier *oidc.IDTokenVerifier
	var err error

	if len(cfg.AllowedIssuers) > 0 {
		// Use the first allowed issuer to initialize the OIDC provider.
		// This assumes that the token endpoint and other discovery URLs
		// are correctly discoverable from this issuer URL.
		// For a truly generic setup, this might need to be more flexible
		// or require the issuer to match the token URL's issuer.
		provider, err = oidc.NewProvider(context.Background(), cfg.AllowedIssuers[0])
		if err != nil {
			// If OIDC provider setup fails, we can still proceed with basic OAuth2
			// but ID token verification will not be possible.
			// Log a warning here.
			fmt.Printf("Warning: Could not initialize OIDC provider for issuer %s: %v. ID token verification will be disabled.\n", cfg.AllowedIssuers[0], err)
		} else {
			oidcConfig := &oidc.Config{
				ClientID: cfg.OAuthClientID,
			}
			verifier = provider.Verifier(oidcConfig)
		}
	} else {
		fmt.Println("Warning: No ALLOWED_ISSUERS configured. ID token verification will be disabled.")
	}

	return &AuthHandler{
		oauth2Config:   oauth2Config,
		provider:       provider,
		verifier:       verifier,
		allowedIssuers: cfg.AllowedIssuers,
		sessionSecret:  cfg.SessionSecret,
	}, nil
}

// LoginHandler redirects the user to the OAuth provider's login page.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState() // Implement generateRandomState
	// TODO: Store state in a short-lived cache or cookie to verify in callback.
	// For simplicity, we are passing it directly, but this is vulnerable to CSRF.
	// A better approach is to store it in an HTTP-only, secure cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   true, // Set to false if testing on HTTP without HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	url := h.oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// CallbackHandler handles the OAuth provider's callback.
func (h *AuthHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Verify state
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, "State cookie not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != stateCookie.Value {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}
	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found in request", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for tokens.
	oauth2Token, err := h.oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Extract the ID Token from the OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "ID token not found in token response", http.StatusInternalServerError)
		return
	}

	var user *User
	if h.verifier != nil {
		// Parse and verify the ID Token payload.
		idToken, err := h.verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to verify ID token: %v", err), http.StatusInternalServerError)
			return
		}

		// Extract custom claims
		if err := idToken.Claims(&user); err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse ID token claims: %v", err), http.StatusInternalServerError)
			return
		}
		// Validate issuer against the allowed list
		issuerValid := false
		for _, allowedIssuer := range h.allowedIssuers {
			if idToken.Issuer == allowedIssuer {
				issuerValid = true
				break
			}
		}
		if !issuerValid {
			http.Error(w, fmt.Sprintf("Invalid issuer: %s", idToken.Issuer), http.StatusUnauthorized)
			return
		}
	} else {
		// Fallback if OIDC verifier is not available (e.g., non-OpenID Connect provider)
		// This is less secure as we don't verify the ID token.
		// We might only have an access token.
		// For now, we'll create a minimal user object.
		// A more robust solution might involve fetching user info from a userinfo endpoint
		// if available, or relying solely on the access token for API calls.
		user = &User{
			ID: "unknown", // Or some other identifier from the access token if possible
		}
		// Log a warning that we are operating in a less secure mode.
		fmt.Println("Warning: ID token verification is disabled. User identity is not fully verified.")
	}

	// Create a session token (JWT) for the user.
	expirationTime := time.Now().Add(1 * time.Hour)
	jwtToken, err := CreateJWT(user, h.sessionSecret, expirationTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create session token: %v", err), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    jwtToken,
		Path:     "/",
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   true, // Set to false if testing on HTTP without HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to the frontend or a success page.
	// This URL should be configurable.
	http.Redirect(w, r, "/", http.StatusFound)
}

// LogoutHandler clears the session cookie and logs the user out.
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

// generateRandomState generates a random string for the OAuth state parameter.
// This is a placeholder. A real implementation should use crypto/rand.
func generateRandomState() string {
	// In a real application, use crypto/rand for a secure random string.
	// For example:
	// b := make([]byte, 16)
	// rand.Read(b)
	// return base64.URLEncoding.EncodeToString(b)
	return "random_state_string_placeholder" // Replace with secure random string
}

// SessionSecret returns the session secret used by the AuthHandler.
// This is needed by the AuthMiddleware.
func (h *AuthHandler) SessionSecret() string {
	return h.sessionSecret
}
