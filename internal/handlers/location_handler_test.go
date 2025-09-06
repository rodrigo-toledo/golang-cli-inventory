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

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocationService is a mock implementation of service.LocationServiceInterface
type MockLocationService struct {
	mock.Mock
}

func (m *MockLocationService) CreateLocation(ctx context.Context, req *models.CreateLocationRequest) (*models.Location, error) {
	args := m.Called(ctx, req)
	// Handle case where location might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationService) GetLocationByName(ctx context.Context, name string) (*models.Location, error) {
	args := m.Called(ctx, name)
	// Handle case where location might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationService) ListLocations(ctx context.Context) ([]models.Location, error) {
	args := m.Called(ctx)
	// Handle case where location list might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Location), args.Error(1)
}

func TestLocationHandler_CreateLocation(t *testing.T) {
	mockService := new(MockLocationService)
	handler := NewLocationHandler(mockService)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreateLocationRequest{
			Name: "Test Location",
		}
		expectedLocation := &models.Location{
			ID:        1,
			Name:      reqBody.Name,
			CreatedAt: time.Now(),
		}

		mockService.On("CreateLocation", mock.Anything, mock.MatchedBy(func(req *models.CreateLocationRequest) bool {
			return req != nil && req.Name == reqBody.Name
		})).Return(expectedLocation, nil)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/locations", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateLocation(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var respLocation models.Location
		err := json.Unmarshal(w.Body.Bytes(), &respLocation)
		assert.NoError(t, err)
		assert.Equal(t, expectedLocation.ID, respLocation.ID)
		assert.Equal(t, expectedLocation.Name, respLocation.Name)
		assert.WithinDuration(t, expectedLocation.CreatedAt, respLocation.CreatedAt, time.Second)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON Payload", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/api/v1/locations", bytes.NewBuffer([]byte("invalid json")))
		w := httptest.NewRecorder()

		handler.CreateLocation(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateLocation")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		reqBody := models.CreateLocationRequest{} // Missing Name
		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/locations", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateLocation(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		resp := w.Body.String()
		assert.Contains(t, resp, "Name is required")
		mockService.AssertNotCalled(t, "CreateLocation")
	})

	t.Run("Service Error", func(t *testing.T) {
		reqBody := models.CreateLocationRequest{
			Name: "Test Location Error",
		}
		mockService.On("CreateLocation", mock.Anything, mock.MatchedBy(func(req *models.CreateLocationRequest) bool {
			return req != nil && req.Name == reqBody.Name
		})).Return((*models.Location)(nil), assert.AnError)

		jsonReq, _ := json.Marshal(reqBody)
		r, _ := http.NewRequest("POST", "/api/v1/locations", bytes.NewBuffer(jsonReq))
		w := httptest.NewRecorder()

		handler.CreateLocation(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLocationHandler_ListLocations(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(MockLocationService)
		handler := NewLocationHandler(mockService)

		expectedLocations := []models.Location{
			{ID: 1, Name: "Location 1", CreatedAt: time.Now()},
			{ID: 2, Name: "Location 2", CreatedAt: time.Now()},
		}
		mockService.On("ListLocations", mock.Anything).Return(expectedLocations, nil)

		r, _ := http.NewRequest("GET", "/api/v1/locations", nil)
		w := httptest.NewRecorder()

		handler.ListLocations(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var respLocations []models.Location
		err := json.Unmarshal(w.Body.Bytes(), &respLocations)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedLocations), len(respLocations))
		for i := range expectedLocations {
			assert.Equal(t, expectedLocations[i].ID, respLocations[i].ID)
			assert.Equal(t, expectedLocations[i].Name, respLocations[i].Name)
			assert.WithinDuration(t, expectedLocations[i].CreatedAt, respLocations[i].CreatedAt, time.Second)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockService := new(MockLocationService)
		handler := NewLocationHandler(mockService)

		mockService.On("ListLocations", mock.Anything).Return(([]models.Location)(nil), assert.AnError)

		r, _ := http.NewRequest("GET", "/api/v1/locations", nil)
		w := httptest.NewRecorder()

		handler.ListLocations(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestLocationHandler_GetLocationByName(t *testing.T) {
	// Setup a minimal chi router for testing URL parameters
	r := chi.NewRouter()
	mockService := new(MockLocationService)
	handler := NewLocationHandler(mockService)
	r.Get("/api/v1/locations/{name}", handler.GetLocationByName)

	t.Run("Success", func(t *testing.T) {
		name := "Test Location"
		expectedLocation := &models.Location{ID: 1, Name: name, CreatedAt: time.Now()}
		mockService.On("GetLocationByName", mock.Anything, name).Return(expectedLocation, nil)

		req, _ := http.NewRequest("GET", "/api/v1/locations/"+name, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respLocation models.Location
		err := json.Unmarshal(w.Body.Bytes(), &respLocation)
		assert.NoError(t, err)
		assert.Equal(t, expectedLocation.ID, respLocation.ID)
		assert.Equal(t, expectedLocation.Name, respLocation.Name)
		assert.WithinDuration(t, expectedLocation.CreatedAt, respLocation.CreatedAt, time.Second)

		mockService.AssertExpectations(t)
	})

	t.Run("Missing Name Param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/locations/", nil)
		w := httptest.NewRecorder()

		handler.GetLocationByName(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Location name is required")
		mockService.AssertNotCalled(t, "GetLocationByName")
	})

	t.Run("Service Error", func(t *testing.T) {
		name := "NonExistent Location"
		mockService.On("GetLocationByName", mock.Anything, name).Return((*models.Location)(nil), assert.AnError)

		req, _ := http.NewRequest("GET", "/api/v1/locations/"+name, nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}