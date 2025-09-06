// Package auth provides authentication and authorization logic for the application.
// It includes OAuth 2.0 flow handling, session management with JWTs,
// and middleware for protecting routes.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT.
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// Authenticator is a middleware that checks for a valid JWT in the request.
func Authenticator(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Try to get the token from the cookie if Authorization header is not present
				cookie, err := r.Cookie("session_token")
				if err != nil || cookie.Value == "" {
					http.Error(w, "Authorization header or session token cookie required", http.StatusUnauthorized)
					return
				}
				authHeader = "Bearer " + cookie.Value
			}

			// The Authorization header should be in the form "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if !(len(parts) == 2 && parts[0] == "Bearer") {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// Parse and validate the JWT
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					http.Error(w, "Token has expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				}
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user information to the request context
			user := &User{
				ID:    claims.UserID,
				Email: claims.Email,
				Name:  claims.Name,
			}
			ctx := context.WithValue(r.Context(), userContextKey, user)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CreateJWT creates a new JWT for the given user.
func CreateJWT(user *User, jwtSecret string, expirationTime time.Time) (string, error) {
	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cli-inventory", // Can be configured
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
