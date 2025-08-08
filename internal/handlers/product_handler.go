// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/go-chi/chi/v5"
)

// ProductHandler handles HTTP requests for product operations.
type ProductHandler struct {
	productService service.ProductServiceInterface
}

// NewProductHandler creates a new instance of ProductHandler.
func NewProductHandler(productService service.ProductServiceInterface) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles POST /api/v1/products requests.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, err) // Will result in a 400 Bad Request
		return
	}

	// TODO: Add more robust validation (e.g., using go-playground/validator)
	if req.SKU == "" || req.Name == "" {
		// For now, we create a simple error to be handled by the generic handler.
		// This can be improved with a specific validation error type.
		HandleError(w, fmt.Errorf("%w: SKU and Name are required", ErrBadRequest))
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		HandleError(w, err) // Handles specific errors like 409 Conflict or 500 Internal Server Error
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		// Log error, but the response header is already sent
		// log.Printf("Failed to encode response: %v", err)
	}
}

// ListProducts handles GET /api/v1/products requests.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.productService.ListProducts(r.Context())
	if err != nil {
		HandleError(w, err) // Handles 500 Internal Server Error
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// GetProductBySKU handles GET /api/v1/products/{sku} requests.
func (h *ProductHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sku := chi.URLParam(r, "sku")
	if sku == "" {
		HandleError(w, fmt.Errorf("%w: SKU is required", ErrBadRequest)) // Will result in a 400 Bad Request
		return
	}

	product, err := h.productService.GetProductBySKU(r.Context(), sku)
	if err != nil {
		HandleError(w, err) // Handles 404 Not Found or 500 Internal Server Error
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}
