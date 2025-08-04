// Package handlers provides HTTP request handlers for the inventory management API.
// It contains handlers for products, locations, and stock operations.
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cli-inventory/internal/models"
	"cli-inventory/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductService is a mock implementation of service.ProductServiceInterface
type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, req)
	// Handle case where product might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	args := m.Called(ctx, sku)
	// Handle case where product might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductService) ListProducts(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	// Handle case where product list might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreateProductRequest{
			SKU:         "TEST-SKU-123",
			Name:        "Test Product",
			Description: "A test product",
			Price:       99.99,
		}
		expectedProduct := &models.Product{
			ID:          1,
			SKU:         reqBody.SKU,
			Name:        reqBody.Name,
			Description: reqBody.Description,
			Price:       reqBody.Price,
		}

		mockService.On("CreateProduct", mock.Anything, mock.MatchedBy(func(req *models.CreateProductRequest) bool {
			return req != nil && req.SKU == reqBody.SKU && req.Name == reqBody.Name
		})).Return(expectedProduct, nil)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var respProduct models.Product
		err := json.NewDecoder(w.Body).Decode(&respProduct)
		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, &respProduct)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateProduct")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		reqBody := models.CreateProductRequest{ // Missing SKU and Name
			Description: "A test product",
			Price:       99.99,
		}
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "SKU and Name are required")
		mockService.AssertNotCalled(t, "CreateProduct")
	})

	t.Run("Service Error", func(t *testing.T) {
		reqBody := models.CreateProductRequest{
			SKU:         "TEST-SKU-ERR",
			Name:        "Test Product Error",
			Description: "A test product for error case",
			Price:       99.99,
		}
		mockService.On("CreateProduct", mock.Anything, mock.MatchedBy(func(req *models.CreateProductRequest) bool {
			return req != nil && req.SKU == reqBody.SKU
		})).Return((*models.Product)(nil), assert.AnError)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_ListProducts(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockProductService)
		handler := NewProductHandler(mockService)

		expectedProducts := []models.Product{
			{ID: 1, SKU: "SKU1", Name: "Product 1", Price: 10.0},
			{ID: 2, SKU: "SKU2", Name: "Product 2", Price: 20.0},
		}
		mockService.On("ListProducts", mock.Anything).Return(expectedProducts, nil)

		r, _ := http.NewRequest("GET", "/api/v1/products", nil)
		w := httptest.NewRecorder()

		handler.ListProducts(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var respProducts []models.Product
		err := json.NewDecoder(w.Body).Decode(&respProducts)
		assert.NoError(t, err)
		assert.Equal(t, expectedProducts, respProducts)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockProductService) // Create a new mock for this sub-test
		handler := NewProductHandler(mockService)

		mockService.On("ListProducts", mock.Anything).Return(([]models.Product)(nil), assert.AnError)

		r, _ := http.NewRequest("GET", "/api/v1/products", nil)
		w := httptest.NewRecorder()

		handler.ListProducts(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected status code 500, got %d. Response body: %s", w.Code, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetProductBySKU(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	// Setup a minimal chi router for testing URL parameters
	r := chi.NewRouter()
	r.Get("/api/v1/products/{sku}", handler.GetProductBySKU)

	t.Run("Success", func(t *testing.T) {
		sku := "TEST-SKU-123"
		expectedProduct := &models.Product{ID: 1, SKU: sku, Name: "Test Product", Price: 99.99}
		mockService.On("GetProductBySKU", mock.Anything, sku).Return(expectedProduct, nil)

		req, _ := http.NewRequest("GET", "/api/v1/products/"+sku, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var respProduct models.Product
		err := json.NewDecoder(w.Body).Decode(&respProduct)
		assert.NoError(t, err)
		assert.Equal(t, *expectedProduct, respProduct)
		mockService.AssertExpectations(t)
	})

	t.Run("Missing SKU Param", func(t *testing.T) {
		// This test case is now implicitly covered by the router's 404 if the path doesn't match,
		// or if the handler itself checks for an empty SKU.
		// The current handler logic checks for empty SKU after chi.URLParam.
		// To test this, we can make a request that doesn't match the route pattern,
		// or make a request that matches but results in an empty SKU.
		// Let's test the handler's direct response to an empty SKU.
		req, _ := http.NewRequest("GET", "/api/v1/products/", nil) // This won't match the route
		w := httptest.NewRecorder()

		// We need a route that can result in an empty SKU param for the handler
		// For this specific test, we call the handler directly with a request that has no SKU param
		// because the router would not even call the handler for "/api/v1/products/"
		// A better way is to have a route like `/api/v1/products/{sku?}` but chi doesn't support optional params easily.
		// So, we test the handler's logic directly for this specific case.
		handler.GetProductBySKU(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "SKU is required")
		mockService.AssertNotCalled(t, "GetProductBySKU")
	})

	t.Run("Service Error - Not Found", func(t *testing.T) {
		sku := "NONEXISTENT-SKU"
		mockService.On("GetProductBySKU", mock.Anything, sku).Return((*models.Product)(nil), service.ErrProductNotFound)

		req, _ := http.NewRequest("GET", "/api/v1/products/"+sku, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// This assertion will change to http.StatusNotFound once the handler is updated
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
