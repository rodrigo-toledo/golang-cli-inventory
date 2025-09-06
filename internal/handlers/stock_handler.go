// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"
)

// StockHandler handles HTTP requests for stock operations.
type StockHandler struct {
	stockService service.StockServiceInterface
}

// NewStockHandler creates a new instance of StockHandler.
func NewStockHandler(stockService service.StockServiceInterface) *StockHandler {
	return &StockHandler{
		stockService: stockService,
	}
}

// AddStock handles POST /api/v1/stock/add requests.
func (h *StockHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req models.AddStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.ProductID <= 0 || req.LocationID <= 0 || req.Quantity <= 0 {
		http.Error(w, "ProductID, LocationID (positive integers) and Quantity (positive integer) are required", http.StatusBadRequest)
		return
	}

	stock, err := h.stockService.AddStock(r.Context(), &req)
	if err != nil {
		// TODO: Handle specific errors (e.g., product/location not found) with appropriate status codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Or 201 Created if we consider this a new stock entry creation
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// MoveStock handles POST /api/v1/stock/move requests.
func (h *StockHandler) MoveStock(w http.ResponseWriter, r *http.Request) {
	var req models.MoveStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.ProductID <= 0 || req.FromLocationID <= 0 || req.ToLocationID <= 0 || req.Quantity <= 0 {
		http.Error(w, "ProductID, FromLocationID, ToLocationID (positive integers) and Quantity (positive integer) are required", http.StatusBadRequest)
		return
	}

	if req.FromLocationID == req.ToLocationID {
		http.Error(w, "Source and destination locations cannot be the same", http.StatusBadRequest)
		return
	}

	stock, err := h.stockService.MoveStock(r.Context(), &req)
	if err != nil {
		// TODO: Handle specific errors (e.g., insufficient stock, product/location not found) with appropriate status codes
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}

// GetLowStockReport handles GET /api/v1/stock/low-stock requests.
func (h *StockHandler) GetLowStockReport(w http.ResponseWriter, r *http.Request) {
	thresholdStr := r.URL.Query().Get("threshold")
	threshold := 10 // Default threshold
	var err error
	if thresholdStr != "" {
		threshold, err = strconv.Atoi(thresholdStr)
		if err != nil || threshold < 0 {
			http.Error(w, "Invalid threshold value, must be a non-negative integer", http.StatusBadRequest)
			return
		}
	}

	stocks, err := h.stockService.GetLowStockReport(r.Context(), threshold)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stocks); err != nil {
		// Log error
		// log.Printf("Failed to encode response: %v", err)
	}
}
