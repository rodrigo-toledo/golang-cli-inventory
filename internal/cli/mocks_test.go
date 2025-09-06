package cli

import (
	"context"
	"cli-inventory/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockProductRepository is a mock implementation of service.ProductRepositoryInterface
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	args := m.Called(ctx, sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

// MockLocationRepository is a mock implementation of service.LocationRepositoryInterface
type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) Create(ctx context.Context, req *models.CreateLocationRequest) (*models.Location, error) {
	args := m.Called(ctx, req)
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

// MockStockRepository is a mock implementation of service.StockRepositoryInterface
type MockStockRepository struct {
	mock.Mock
}

func (m *MockStockRepository) AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	args := m.Called(ctx, productID, locationID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Stock), args.Error(1)
}

func (m *MockStockRepository) RemoveStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	args := m.Called(ctx, productID, locationID, quantity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetLowStock(ctx context.Context, threshold int) ([]models.Stock, error) {
	args := m.Called(ctx, threshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Stock), args.Error(1)
}

func (m *MockStockRepository) GetByProductAndLocation(ctx context.Context, productID, locationID int) (*models.Stock, error) {
	args := m.Called(ctx, productID, locationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Stock), args.Error(1)
}

// MockStockMovementRepository is a mock implementation of service.StockMovementRepositoryInterface
type MockStockMovementRepository struct {
	mock.Mock
}

func (m *MockStockMovementRepository) Create(ctx context.Context, movement *models.StockMovement) (*models.StockMovement, error) {
	args := m.Called(ctx, movement)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.StockMovement), args.Error(1)
}