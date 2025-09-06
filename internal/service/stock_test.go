package service

import (
	"context"
	"testing"

	"github.com/rodrigotoledo/cli-inventory/internal/models"
)

// MockStockProductRepository is a mock implementation of ProductRepositoryInterface for testing
type MockStockProductRepository struct {
	products map[int]*models.Product
}

func (m *MockStockProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if p, exists := m.products[id]; exists {
		return p, nil
	}
	return nil, nil // Simulate not found
}

func (m *MockStockProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
	// This is a simplified mock implementation
	return nil, nil
}

func (m *MockStockProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	// This is a simplified mock implementation
	return nil, nil
}

func (m *MockStockProductRepository) List(ctx context.Context) ([]models.Product, error) {
	// This is a simplified mock implementation
	return nil, nil
}

// MockStockLocationRepository is a mock implementation of LocationRepository for testing
type MockStockLocationRepository struct {
	locations map[int]*models.Location
}

func (m *MockStockLocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	if l, exists := m.locations[id]; exists {
		return l, nil
	}
	return nil, nil // Simulate not found
}

// MockStockRepositoryImpl is a mock implementation of StockRepository for testing
type MockStockRepositoryImpl struct {
	stock map[[2]int]*models.Stock // key: [productID, locationID]
}

func (m *MockStockRepositoryImpl) AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	key := [2]int{productID, locationID}
	if s, exists := m.stock[key]; exists {
		s.Quantity += quantity
		return s, nil
	}

	s := &models.Stock{
		ID:         len(m.stock) + 1,
		ProductID:  productID,
		LocationID: locationID,
		Quantity:   quantity,
	}
	m.stock[key] = s
	return s, nil
}

func (m *MockStockRepositoryImpl) RemoveStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	key := [2]int{productID, locationID}
	if s, exists := m.stock[key]; exists {
		s.Quantity -= quantity
		if s.Quantity < 0 {
			s.Quantity = 0
		}
		return s, nil
	}

	// If stock doesn't exist, create it with 0 quantity
	s := &models.Stock{
		ID:         len(m.stock) + 1,
		ProductID:  productID,
		LocationID: locationID,
		Quantity:   0,
	}
	m.stock[key] = s
	return s, nil
}

func (m *MockStockRepositoryImpl) GetLowStock(ctx context.Context, threshold int) ([]models.Stock, error) {
	stocks := make([]models.Stock, 0)
	for _, s := range m.stock {
		if s.Quantity < threshold {
			stocks = append(stocks, *s)
		}
	}
	return stocks, nil
}

// MockStockMovementRepositoryImpl is a mock implementation of StockMovementRepository for testing
type MockStockMovementRepositoryImpl struct {
	movements []models.StockMovement
}

func (m *MockStockMovementRepositoryImpl) Create(ctx context.Context, movement *models.StockMovement) (*models.StockMovement, error) {
	movement.ID = len(m.movements) + 1
	m.movements = append(m.movements, *movement)
	return movement, nil
}

func TestStockService_AddStock(t *testing.T) {
	productRepo := &MockStockProductRepository{
		products: map[int]*models.Product{
			1: {ID: 1, SKU: "TEST001", Name: "Test Product"},
		},
	}

	locationRepo := &MockStockLocationRepository{
		locations: map[int]*models.Location{
			1: {ID: 1, Name: "Test Location"},
		},
	}

	stockRepo := &MockStockRepositoryImpl{
		stock: make(map[[2]int]*models.Stock),
	}

	movementRepo := &MockStockMovementRepositoryImpl{
		movements: make([]models.StockMovement, 0),
	}

	// For this test, we'll pass nil for the db parameter since we're not using it
	service := NewStockService(productRepo, locationRepo, stockRepo, movementRepo, nil)

	ctx := context.Background()
	req := &models.AddStockRequest{
		ProductID:  1,
		LocationID: 1,
		Quantity:   10,
	}

	stock, err := service.AddStock(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if stock.Quantity != 10 {
		t.Errorf("Expected quantity 10, got %d", stock.Quantity)
	}

	if stock.ProductID != 1 {
		t.Errorf("Expected ProductID 1, got %d", stock.ProductID)
	}

	if stock.LocationID != 1 {
		t.Errorf("Expected LocationID 1, got %d", stock.LocationID)
	}
}

func TestStockService_MoveStock(t *testing.T) {
	productRepo := &MockStockProductRepository{
		products: map[int]*models.Product{
			1: {ID: 1, SKU: "TEST001", Name: "Test Product"},
		},
	}

	locationRepo := &MockStockLocationRepository{
		locations: map[int]*models.Location{
			1: {ID: 1, Name: "Source Location"},
			2: {ID: 2, Name: "Destination Location"},
		},
	}

	stockRepo := &MockStockRepositoryImpl{
		stock: map[[2]int]*models.Stock{
			[2]int{1, 1}: {ID: 1, ProductID: 1, LocationID: 1, Quantity: 10},
		},
	}

	movementRepo := &MockStockMovementRepositoryImpl{
		movements: make([]models.StockMovement, 0),
	}

	// For this test, we'll pass nil for the db parameter since we're not using it
	service := NewStockService(productRepo, locationRepo, stockRepo, movementRepo, nil)

	ctx := context.Background()
	req := &models.MoveStockRequest{
		ProductID:      1,
		FromLocationID: 1,
		ToLocationID:   2,
		Quantity:       5,
	}

	stock, err := service.MoveStock(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the stock was moved
	if stock.Quantity != 5 {
		t.Errorf("Expected quantity 5 at destination, got %d", stock.Quantity)
	}
}
