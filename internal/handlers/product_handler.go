// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json"
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
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: Add more robust validation (e.g., using go-playground/validator)
	if req.SKU == "" || req.Name == "" {
		http.Error(w, "SKU and Name are required", http.StatusBadRequest)
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		// TODO: Handle specific errors (e.g., product already exists) with appropriate status codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		// Log error, but the response header is already sent
		// log.Printf("Failed to encode response: %v", err)
	}
}

// ListProducts handles GET /api/v1/products requests.
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// GetProductBySKU handles GET /api/v1/products/{sku} requests.
func (h *ProductHandler) GetProductBySKU(w http.ResponseWriter, r *http.Request) {
	sku := chi.URLParam(r, "sku")
	if sku == "" {
		http.Error(w, "SKU is required", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProductBySKU(r.Context(), sku)
	if err != nil {
		// TODO: Check for "not found" error and return 404
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}
