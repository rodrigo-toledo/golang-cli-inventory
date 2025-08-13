// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json/v2"
	"net/http"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/go-chi/chi/v5"
)

// LocationHandler handles HTTP requests for location operations.
type LocationHandler struct {
	locationService service.LocationServiceInterface
}

// NewLocationHandler creates a new instance of LocationHandler.
func NewLocationHandler(locationService service.LocationServiceInterface) *LocationHandler {
	return &LocationHandler{
		locationService: locationService,
	}
}

// CreateLocation handles POST /api/v1/locations requests.
func (h *LocationHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLocationRequest
	if err := json.UnmarshalRead(r.Body, &req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	location, err := h.locationService.CreateLocation(r.Context(), &req)
	if err != nil {
		// TODO: Handle specific errors (e.g., location already exists) with appropriate status codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.MarshalWrite(w, location); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// ListLocations handles GET /api/v1/locations requests.
func (h *LocationHandler) ListLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.locationService.ListLocations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.MarshalWrite(w, locations); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// GetLocationByName handles GET /api/v1/locations/{name} requests.
func (h *LocationHandler) GetLocationByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, "Location name is required", http.StatusBadRequest)
		return
	}

	location, err := h.locationService.GetLocationByName(r.Context(), name)
	if err != nil {
		// TODO: Check for "not found" error and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.MarshalWrite(w, location); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}
