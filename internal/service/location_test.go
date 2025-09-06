package service

import (
	"context"
	"fmt"
	"testing"

	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLocationRepository is a mock implementation that mimics the LocationRepository methods
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Create(ctx context.Context, location *models.CreateLocationRequest) (*models.Location, error) {
	args := m.Called(ctx, location)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) GetByName(ctx context.Context, name string) (*models.Location, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Location), args.Error(1)
}

func (m *MockLocationRepository) List(ctx context.Context) ([]models.Location, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Location), args.Error(1)
}

func TestNewLocationService(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestLocationService_CreateLocation(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := NewLocationService(mockRepo)

	ctx := context.Background()
	req := &models.CreateLocationRequest{
		Name: "Test Location",
	}

	expectedLocation := &models.Location{
		ID:   1,
		Name: req.Name,
	}

	// Test successful creation
	mockRepo.On("GetByName", ctx, req.Name).Return(nil, fmt.Errorf("not found"))
	mockRepo.On("Create", ctx, req).Return(expectedLocation, nil)

	location, err := service.CreateLocation(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, location)

	mockRepo.AssertExpectations(t)
}

func TestLocationService_CreateLocation_AlreadyExists(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	req := &models.CreateLocationRequest{
		Name: "Test Location",
	}

	existingLocation := &models.Location{
		ID:   1,
		Name: req.Name,
	}

	// Test location already exists
	mockRepo.On("GetByName", ctx, req.Name).Return(existingLocation, nil)

	location, err := service.CreateLocation(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "already exists")

	mockRepo.AssertExpectations(t)
}

func TestLocationService_CreateLocation_ErrorOnGetByName(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	req := &models.CreateLocationRequest{
		Name: "Test Location",
	}

	// Test error when checking if location exists, but still try to create
	mockRepo.On("GetByName", ctx, req.Name).Return(nil, fmt.Errorf("database error"))
	mockRepo.On("Create", ctx, req).Return(&models.Location{}, nil)

	location, err := service.CreateLocation(ctx, req)
	// In this case, we expect the Create to be called and succeed
	assert.NoError(t, err)
	assert.NotNil(t, location)

	mockRepo.AssertExpectations(t)
}

func TestLocationService_CreateLocation_ErrorOnCreate(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	req := &models.CreateLocationRequest{
		Name: "Test Location",
	}

	// Test error when creating location
	mockRepo.On("GetByName", ctx, req.Name).Return(nil, fmt.Errorf("not found"))
	mockRepo.On("Create", ctx, req).Return(nil, fmt.Errorf("database error"))

	location, err := service.CreateLocation(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "failed to create location")

	mockRepo.AssertExpectations(t)
}

func TestLocationService_GetLocationByName(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	name := "Test Location"

	expectedLocation := &models.Location{
		ID:   1,
		Name: name,
	}

	// Test successful retrieval
	mockRepo.On("GetByName", ctx, name).Return(expectedLocation, nil)

	location, err := service.GetLocationByName(ctx, name)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocation, location)

	mockRepo.AssertExpectations(t)
}

func TestLocationService_GetLocationByName_NotFound(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	name := "Non-existent Location"

	// Test location not found
	mockRepo.On("GetByName", ctx, name).Return(nil, pgx.ErrNoRows)

	location, err := service.GetLocationByName(ctx, name)
	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "location not found")

	mockRepo.AssertExpectations(t)
}

func TestLocationService_GetLocationByName_Error(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()
	name := "Test Location"

	// Test database error
	mockRepo.On("GetByName", ctx, name).Return(nil, fmt.Errorf("database error"))

	location, err := service.GetLocationByName(ctx, name)
	assert.Error(t, err)
	assert.Nil(t, location)
	assert.Contains(t, err.Error(), "failed to get location")

	mockRepo.AssertExpectations(t)
}

func TestLocationService_ListLocations(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()

	expectedLocations := []models.Location{
		{
			ID:   1,
			Name: "Location 1",
		},
		{
			ID:   2,
			Name: "Location 2",
		},
	}

	// Test successful listing
	mockRepo.On("List", ctx).Return(expectedLocations, nil)

	locations, err := service.ListLocations(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedLocations, locations)

	mockRepo.AssertExpectations(t)
}

func TestLocationService_ListLocations_Error(t *testing.T) {
	mockRepo := new(MockLocationRepository)
	service := &LocationService{repo: mockRepo}

	ctx := context.Background()

	// Test database error
	mockRepo.On("List", ctx).Return(nil, fmt.Errorf("database error"))

	locations, err := service.ListLocations(ctx)
	assert.Error(t, err)
	assert.Nil(t, locations)
	assert.Contains(t, err.Error(), "failed to list locations")

	mockRepo.AssertExpectations(t)
}