// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"strings"

	"cli-inventory/internal/service"
)

// ErrorResponse defines the structure for error responses sent to the client.
// This aligns with the OpenAPI specification's Error schema.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// ErrBadRequest is a generic error for client-side bad requests, e.g., validation failures.
var ErrBadRequest = errors.New("bad request")

// HandleError maps service-level errors to appropriate HTTP status codes and responses.
// It centralizes error response logic to ensure consistency across all handlers.
func HandleError(w http.ResponseWriter, err error) {
	// Check for specific, known errors and map them to HTTP status codes.
	// This list should be expanded as new custom errors are defined in the service layer.
	switch {
	case errors.Is(err, service.ErrProductNotFound):
		respondWithError(w, http.StatusNotFound, "Resource not found", err.Error())
	case errors.Is(err, service.ErrLocationNotFound):
		respondWithError(w, http.StatusNotFound, "Resource not found", err.Error())
	case errors.Is(err, service.ErrInsufficientStock):
		respondWithError(w, http.StatusConflict, "Insufficient stock", err.Error())
	case errors.Is(err, ErrBadRequest):
		// We expect the error to be wrapped with a specific message.
		// e.g. fmt.Errorf("%w: SKU and Name are required", ErrBadRequest)
		respondWithError(w, http.StatusBadRequest, "Invalid request", err.Error())
	case isJSONError(err):
		respondWithError(w, http.StatusBadRequest, "Invalid request payload", err.Error())
	// Example for a potential future error:
	// case errors.Is(err, service.ErrSKUAlreadyExists):
	// 	respondWithError(w, http.StatusConflict, "Product with this SKU already exists", err.Error())
	default:
		// For any other unhandled errors, return a generic 500 Internal Server Error.
		// This prevents leaking sensitive internal error details to the client.
		respondWithError(w, http.StatusInternalServerError, "An internal server error occurred", "Please try again later.")
	}
}

// isJSONError checks if the error is related to JSON decoding.
func isJSONError(err error) bool {
	var syntaxError *json.SemanticError
	var unmarshalTypeError *json.SemanticError
	// Check if it's a syntactic error by checking the error message
	// This is a workaround since we can't directly import jsontext.SyntacticError
	return errors.As(err, &syntaxError) || errors.As(err, &unmarshalTypeError) || 
		(err != nil && (strings.Contains(err.Error(), "jsontext:") || strings.Contains(err.Error(), "invalid character")))
}

// respondWithError is a helper function to send a JSON error response.
func respondWithError(w http.ResponseWriter, code int, message string, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Error:   message,
		Details: details,
	}

	// Use JSON v2 MarshalWrite function
	if err := json.MarshalWrite(w, errorResponse); err != nil {
		// Log the error that occurred while trying to send the error response.
		// This is a fallback, as we can't do much more if encoding fails here.
		// In a real application, a proper logger would be used.
		// log.Printf("Failed to encode error response: %v", err)
	}
}
