// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cli-inventory/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStockService is a mock implementation of service.StockServiceInterface
type MockStockService struct {
	mock.Mock
}

func (m *MockStockService) AddStock(ctx context.Context, req *models.AddStockRequest) (*models.Stock, error) {
	args := m.Called(ctx, req)
	// Handle case where stock might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Stock), args.Error(1)
}

func (m *MockStockService) MoveStock(ctx context.Context, req *models.MoveStockRequest) (*models.Stock, error) {
	args := m.Called(ctx, req)
	// Handle case where stock might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Stock), args.Error(1)
}

func (m *MockStockService) GetLowStockReport(ctx context.Context, threshold int) ([]models.Stock, error) {
	args := m.Called(ctx, threshold)
	// Handle case where stock list might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Stock), args.Error(1)
}

func TestStockHandler_AddStock(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.AddStockRequest{
			ProductID:  1,
			LocationID: 1,
			Quantity:   100,
		}
		expectedStock := &models.Stock{
			ID:         1,
			ProductID:  reqBody.ProductID,
			LocationID: reqBody.LocationID,
			Quantity:   reqBody.Quantity,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockService.On("AddStock", mock.Anything, mock.MatchedBy(func(req *models.AddStockRequest) bool {
			return req != nil && req.ProductID == reqBody.ProductID &&
				req.LocationID == reqBody.LocationID && req.Quantity == reqBody.Quantity
		})).Return(expectedStock, nil)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var respStock models.Stock
		err := json.Unmarshal(w.Body.Bytes(), &respStock)
		assert.NoError(t, err)
		assert.Equal(t, expectedStock.ID, respStock.ID)
		assert.Equal(t, expectedStock.ProductID, respStock.ProductID)
		assert.Equal(t, expectedStock.LocationID, respStock.LocationID)
		assert.Equal(t, expectedStock.Quantity, respStock.Quantity)
		assert.WithinDuration(t, expectedStock.CreatedAt, respStock.CreatedAt, time.Second)
		assert.WithinDuration(t, expectedStock.UpdatedAt, respStock.UpdatedAt, time.Second)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "AddStock")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.AddStockRequest{} // Missing all required fields
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "ProductID, LocationID (positive integers) and Quantity (positive integer) are required")
		mockService.AssertNotCalled(t, "AddStock")
	})

	t.Run("Zero Values", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.AddStockRequest{
			ProductID:  0,
			LocationID: 0,
			Quantity:   0,
		}
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "ProductID, LocationID (positive integers) and Quantity (positive integer) are required")
		mockService.AssertNotCalled(t, "AddStock")
	})

	t.Run("Negative Values", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.AddStockRequest{
			ProductID:  -1,
			LocationID: -1,
			Quantity:   -5,
		}
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "ProductID, LocationID (positive integers) and Quantity (positive integer) are required")
		mockService.AssertNotCalled(t, "AddStock")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.AddStockRequest{
			ProductID:  1,
			LocationID: 1,
			Quantity:   100,
		}
		mockService.On("AddStock", mock.Anything, mock.MatchedBy(func(req *models.AddStockRequest) bool {
			return req != nil && req.ProductID == reqBody.ProductID &&
				req.LocationID == reqBody.LocationID && req.Quantity == reqBody.Quantity
		})).Return((*models.Stock)(nil), assert.AnError)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/add", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.AddStock(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestStockHandler_MoveStock(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.MoveStockRequest{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   2,
			Quantity:       50,
		}
		expectedStock := &models.Stock{
			ID:         1,
			ProductID:  reqBody.ProductID,
			LocationID: reqBody.ToLocationID,
			Quantity:   reqBody.Quantity,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		mockService.On("MoveStock", mock.Anything, mock.MatchedBy(func(req *models.MoveStockRequest) bool {
			return req != nil && req.ProductID == reqBody.ProductID &&
				req.FromLocationID == reqBody.FromLocationID &&
				req.ToLocationID == reqBody.ToLocationID &&
				req.Quantity == reqBody.Quantity
		})).Return(expectedStock, nil)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/move", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.MoveStock(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var respStock models.Stock
		err := json.Unmarshal(w.Body.Bytes(), &respStock)
		assert.NoError(t, err)
		assert.Equal(t, expectedStock.ID, respStock.ID)
		assert.Equal(t, expectedStock.ProductID, respStock.ProductID)
		assert.Equal(t, expectedStock.LocationID, respStock.LocationID)
		assert.Equal(t, expectedStock.Quantity, respStock.Quantity)
		assert.WithinDuration(t, expectedStock.CreatedAt, respStock.CreatedAt, time.Second)
		assert.WithinDuration(t, expectedStock.UpdatedAt, respStock.UpdatedAt, time.Second)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		r, _ := http.NewRequest("POST", "/api/v1/stock/move", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.MoveStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "MoveStock")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.MoveStockRequest{} // Missing all required fields
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/move", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.MoveStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "ProductID, FromLocationID, ToLocationID (positive integers) and Quantity (positive integer) are required")
		mockService.AssertNotCalled(t, "MoveStock")
	})

	t.Run("Same Source and Destination", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.MoveStockRequest{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   1, // Same as FromLocationID
			Quantity:       50,
		}
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/move", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.MoveStock(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "Source and destination locations cannot be the same")
		mockService.AssertNotCalled(t, "MoveStock")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		reqBody := models.MoveStockRequest{
			ProductID:      1,
			FromLocationID: 1,
			ToLocationID:   2,
			Quantity:       50,
		}
		mockService.On("MoveStock", mock.Anything, mock.MatchedBy(func(req *models.MoveStockRequest) bool {
			return req != nil && req.ProductID == reqBody.ProductID &&
				req.FromLocationID == reqBody.FromLocationID &&
				req.ToLocationID == reqBody.ToLocationID &&
				req.Quantity == reqBody.Quantity
		})).Return((*models.Stock)(nil), assert.AnError)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/stock/move", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.MoveStock(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestStockHandler_GetLowStockReport(t *testing.T) {
	t.Run("Success with Default Threshold", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		expectedStocks := []models.Stock{
			{ID: 1, ProductID: 1, LocationID: 1, Quantity: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, ProductID: 2, LocationID: 1, Quantity: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		threshold := 10 // Default threshold

		mockService.On("GetLowStockReport", mock.Anything, threshold).Return(expectedStocks, nil)

		r, _ := http.NewRequest("GET", "/api/v1/stock/low-stock", nil)
		w := httptest.NewRecorder()

		handler.GetLowStockReport(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var respStocks []models.Stock
		err := json.Unmarshal(w.Body.Bytes(), &respStocks)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedStocks), len(respStocks))
		for i := range expectedStocks {
			assert.Equal(t, expectedStocks[i].ID, respStocks[i].ID)
			assert.Equal(t, expectedStocks[i].ProductID, respStocks[i].ProductID)
			assert.Equal(t, expectedStocks[i].LocationID, respStocks[i].LocationID)
			assert.Equal(t, expectedStocks[i].Quantity, respStocks[i].Quantity)
			assert.WithinDuration(t, expectedStocks[i].CreatedAt, respStocks[i].CreatedAt, time.Second)
			assert.WithinDuration(t, expectedStocks[i].UpdatedAt, respStocks[i].UpdatedAt, time.Second)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Success with Custom Threshold", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		expectedStocks := []models.Stock{
			{ID: 1, ProductID: 1, LocationID: 1, Quantity: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		threshold := 20

		mockService.On("GetLowStockReport", mock.Anything, threshold).Return(expectedStocks, nil)

		r, _ := http.NewRequest("GET", "/api/v1/stock/low-stock?threshold=20", nil)
		w := httptest.NewRecorder()

		handler.GetLowStockReport(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var respStocks []models.Stock
		err := json.Unmarshal(w.Body.Bytes(), &respStocks)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedStocks), len(respStocks))

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid Threshold", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		r, _ := http.NewRequest("GET", "/api/v1/stock/low-stock?threshold=invalid", nil)
		w := httptest.NewRecorder()

		handler.GetLowStockReport(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "Invalid threshold value, must be a non-negative integer")
		mockService.AssertNotCalled(t, "GetLowStockReport")
	})

	t.Run("Negative Threshold", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		r, _ := http.NewRequest("GET", "/api/v1/stock/low-stock?threshold=-5", nil)
		w := httptest.NewRecorder()

		handler.GetLowStockReport(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "Invalid threshold value, must be a non-negative integer")
		mockService.AssertNotCalled(t, "GetLowStockReport")
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockStockService)
		handler := NewStockHandler(mockService)

		threshold := 10
		mockService.On("GetLowStockReport", mock.Anything, threshold).Return(([]models.Stock)(nil), assert.AnError)

		r, _ := http.NewRequest("GET", "/api/v1/stock/low-stock", nil)
		w := httptest.NewRecorder()

		handler.GetLowStockReport(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}