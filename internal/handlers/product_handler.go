// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json/v2"
	"fmt"
	"net/http"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/go-chi/chi/v5"
	validator "github.com/go-playground/validator/v10"
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

var validate = validator.New()

// CreateProduct handles POST /api/v1/products requests.
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateProductRequest
	if err := json.UnmarshalRead(r.Body, &req); err != nil {
		HandleError(w, err) // Will result in a 400 Bad Request
		return
	}

	// Validate request using go-playground/validator tags on the model.
	if err := validate.Struct(req); err != nil {
		HandleError(w, fmt.Errorf("%w: %v", ErrBadRequest, err.Error()))
		return
	}

	product, err := h.productService.CreateProduct(r.Context(), &req)
	if err != nil {
		HandleError(w, err) // Handles specific errors like 409 Conflict or 500 Internal Server Error
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.MarshalWrite(w, product); err != nil {
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
	if err := json.MarshalWrite(w, products); err != nil {
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
	if err := json.MarshalWrite(w, product); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}
