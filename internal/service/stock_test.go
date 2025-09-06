package service

import (
	"context"
	"fmt"
	"testing"

	"cli-inventory/internal/models"
)

// MockStockProductRepository is a mock implementation of ProductRepositoryInterface for testing
type MockStockProductRepository struct {
	products map[int]*models.Product
}

func (m *MockStockProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if p, exists := m.products[id]; exists {
		return p, nil
	}
	return nil, fmt.Errorf("product with ID %d not found", id)
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

// MockStockLocationRepository is a mock implementation of LocationRepositoryInterface for testing
type MockStockLocationRepository struct {
	locations map[int]*models.Location
}

func (m *MockStockLocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	if l, exists := m.locations[id]; exists {
		return l, nil
	}
	return nil, fmt.Errorf("location with ID %d not found", id)
}

func (m *MockStockLocationRepository) Create(ctx context.Context, location *models.CreateLocationRequest) (*models.Location, error) {
	// This method is not used in stock tests, so we can return a basic implementation
	id := len(m.locations) + 1
	l := &models.Location{
		ID:   id,
		Name: location.Name,
	}
	if m.locations == nil {
		m.locations = make(map[int]*models.Location)
	}
	m.locations[id] = l
	return l, nil
}

func (m *MockStockLocationRepository) GetByName(ctx context.Context, name string) (*models.Location, error) {
	// This method is not used in stock tests, so we can return a basic implementation
	return nil, fmt.Errorf("location with name %s not found", name)
}

func (m *MockStockLocationRepository) List(ctx context.Context) ([]models.Location, error) {
	// This method is not used in stock tests, so we can return a basic implementation
	return []models.Location{}, nil
}

// MockStockRepositoryImpl is a mock implementation of StockRepository for testing
type MockStockRepositoryImpl struct {
	stock map[[2]int]*models.Stock // key: [productID, locationID]
}

func (m *MockStockRepositoryImpl) AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

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
		// For negative thresholds, all stock items should be considered low
		if threshold < 0 || s.Quantity < threshold {
			stocks = append(stocks, *s)
		}
	}
	return stocks, nil
}

func (m *MockStockRepositoryImpl) GetByProductAndLocation(ctx context.Context, productID, locationID int) (*models.Stock, error) {
	key := [2]int{productID, locationID}
	if s, exists := m.stock[key]; exists {
		return s, nil
	}
	return nil, fmt.Errorf("stock not found for product %d at location %d", productID, locationID)
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

func TestStockService_AddStock_InvalidInput(t *testing.T) {
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

	service := NewStockService(productRepo, locationRepo, stockRepo, movementRepo, nil)

	ctx := context.Background()

	testCases := []struct {
		name    string
		req     *models.AddStockRequest
		wantErr bool
	}{
		{
			name: "Non-existent Product",
			req: &models.AddStockRequest{
				ProductID:  999,
				LocationID: 1,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Non-existent Location",
			req: &models.AddStockRequest{
				ProductID:  1,
				LocationID: 999,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Negative Quantity",
			req: &models.AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   -5,
			},
			wantErr: true,
		},
		{
			name: "Zero Quantity",
			req: &models.AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   0,
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.AddStock(ctx, tc.req)

			if tc.wantErr && err == nil {
				t.Fatalf("Expected error, got none")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestStockService_MoveStock_InvalidInput(t *testing.T) {
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

	service := NewStockService(productRepo, locationRepo, stockRepo, movementRepo, nil)

	ctx := context.Background()

	testCases := []struct {
		name    string
		req     *models.MoveStockRequest
		wantErr bool
	}{
		{
			name: "Non-existent Product",
			req: &models.MoveStockRequest{
				ProductID:      999,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Non-existent Source Location",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 999,
				ToLocationID:   2,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Non-existent Destination Location",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   999,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Same Source and Destination",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   1,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Negative Quantity",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       -5,
			},
			wantErr: true,
		},
		{
			name: "Zero Quantity",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       0,
			},
			wantErr: true,
		},
		{
			name: "Insufficient Stock",
			req: &models.MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       20, // More than available (10)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.MoveStock(ctx, tc.req)

			if tc.wantErr && err == nil {
				t.Fatalf("Expected error, got none")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}

func TestStockService_GetLowStockReport(t *testing.T) {
	productRepo := &MockStockProductRepository{
		products: map[int]*models.Product{
			1: {ID: 1, SKU: "LOW1", Name: "Low Stock Product 1"},
			2: {ID: 2, SKU: "LOW2", Name: "Low Stock Product 2"},
			3: {ID: 3, SKU: "HIGH1", Name: "High Stock Product 1"},
		},
	}

	locationRepo := &MockStockLocationRepository{
		locations: map[int]*models.Location{
			1: {ID: 1, Name: "Warehouse A"},
			2: {ID: 2, Name: "Warehouse B"},
		},
	}

	stockRepo := &MockStockRepositoryImpl{
		stock: map[[2]int]*models.Stock{
			[2]int{1, 1}: {ID: 1, ProductID: 1, LocationID: 1, Quantity: 5},  // Low stock
			[2]int{2, 1}: {ID: 2, ProductID: 2, LocationID: 1, Quantity: 8},  // Low stock
			[2]int{3, 1}: {ID: 3, ProductID: 3, LocationID: 1, Quantity: 50}, // High stock
			[2]int{1, 2}: {ID: 4, ProductID: 1, LocationID: 2, Quantity: 15}, // High stock
		},
	}

	movementRepo := &MockStockMovementRepositoryImpl{
		movements: make([]models.StockMovement, 0),
	}

	service := NewStockService(productRepo, locationRepo, stockRepo, movementRepo, nil)

	ctx := context.Background()

	t.Run("Get Low Stock Report", func(t *testing.T) {
		threshold := 10
		lowStock, err := service.GetLowStockReport(ctx, threshold)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(lowStock) != 2 {
			t.Fatalf("Expected 2 low stock items, got %d", len(lowStock))
		}

		// Verify the correct items are returned
		stockMap := make(map[[2]int]int)
		for _, s := range lowStock {
			stockMap[[2]int{s.ProductID, s.LocationID}] = s.Quantity
		}

		if stockMap[[2]int{1, 1}] != 5 {
			t.Errorf("Expected quantity 5 for product 1 at location 1, got %d", stockMap[[2]int{1, 1}])
		}

		if stockMap[[2]int{2, 1}] != 8 {
			t.Errorf("Expected quantity 8 for product 2 at location 1, got %d", stockMap[[2]int{2, 1}])
		}
	})

	t.Run("Get Low Stock Report With High Threshold", func(t *testing.T) {
		threshold := 100
		lowStock, err := service.GetLowStockReport(ctx, threshold)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// All items should be considered low stock with threshold of 100
		if len(lowStock) != 4 {
			t.Fatalf("Expected 4 low stock items with high threshold, got %d", len(lowStock))
		}
	})

	t.Run("Get Low Stock Report With Zero Threshold", func(t *testing.T) {
		threshold := 0
		lowStock, err := service.GetLowStockReport(ctx, threshold)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// No items should be considered low stock with threshold of 0
		if len(lowStock) != 0 {
			t.Fatalf("Expected 0 low stock items with zero threshold, got %d", len(lowStock))
		}
	})

	t.Run("Get Low Stock Report With Negative Threshold", func(t *testing.T) {
		threshold := -5
		lowStock, err := service.GetLowStockReport(ctx, threshold)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// All items should be considered low stock with negative threshold
		if len(lowStock) != 4 {
			t.Fatalf("Expected 4 low stock items with negative threshold, got %d", len(lowStock))
		}
	})
}
