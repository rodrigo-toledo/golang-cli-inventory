package service

import (
	"context"
	"fmt"
	"testing"

	"cli-inventory/internal/models"
)

// MockProductRepository is a mock implementation of ProductRepositoryInterface for testing
type MockProductRepository struct {
	products map[string]*models.Product
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
	// Basic validation
	if product.SKU == "" {
		return nil, fmt.Errorf("SKU cannot be empty")
	}
	if product.Name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if product.Price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}

	// Check for duplicate SKU
	if _, exists := m.products[product.SKU]; exists {
		return nil, fmt.Errorf("product with SKU %s already exists", product.SKU)
	}

	p := &models.Product{
		ID:          len(m.products) + 1,
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	m.products[product.SKU] = p
	return p, nil
}

func (m *MockProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	if p, exists := m.products[sku]; exists {
		return p, nil
	}
	return nil, nil // Simulate not found
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	for _, p := range m.products {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, nil // Simulate not found
}

func (m *MockProductRepository) List(ctx context.Context) ([]models.Product, error) {
	products := make([]models.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, *p)
	}
	return products, nil
}

func TestProductService_CreateProduct(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()
	req := &models.CreateProductRequest{
		SKU:         "TEST001",
		Name:        "Test Product",
		Description: "A test product",
		Price:       9.99,
	}

	product, err := service.CreateProduct(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if product.SKU != req.SKU {
		t.Errorf("Expected SKU %s, got %s", req.SKU, product.SKU)
	}

	if product.Name != req.Name {
		t.Errorf("Expected Name %s, got %s", req.Name, product.Name)
	}
}

func TestProductService_CreateProduct_DuplicateSKU(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()
	req := &models.CreateProductRequest{
		SKU:         "TEST001",
		Name:        "Test Product",
		Description: "A test product",
		Price:       9.99,
	}

	// Create the product first
	_, err := service.CreateProduct(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error on first create, got %v", err)
	}

	// Try to create the same product again
	_, err = service.CreateProduct(ctx, req)
	if err == nil {
		t.Fatalf("Expected error for duplicate SKU, got none")
	}
}

func TestProductService_GetProductBySKU(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()
	req := &models.CreateProductRequest{
		SKU:         "TEST001",
		Name:        "Test Product",
		Description: "A test product",
		Price:       9.99,
	}

	// Create a product
	createdProduct, err := service.CreateProduct(ctx, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Retrieve the product
	product, err := service.GetProductBySKU(ctx, "TEST001")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if product.SKU != createdProduct.SKU {
		t.Errorf("Expected SKU %s, got %s", createdProduct.SKU, product.SKU)
	}

	if product.Name != createdProduct.Name {
		t.Errorf("Expected Name %s, got %s", createdProduct.Name, product.Name)
	}
}

func TestProductService_GetProductBySKU_NotFound(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()

	// Try to retrieve a product that doesn't exist
	product, err := service.GetProductBySKU(ctx, "NONEXISTENT")
	if err != nil {
		t.Fatalf("Expected no error for not found, got %v", err)
	}

	if product != nil {
		t.Fatalf("Expected nil product for not found, got %v", product)
	}
}

func TestProductService_ListProducts(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()

	// Create multiple products
	products := []*models.CreateProductRequest{
		{
			SKU:         "LIST001",
			Name:        "Product 1",
			Description: "First product",
			Price:       5.99,
		},
		{
			SKU:         "LIST002",
			Name:        "Product 2",
			Description: "Second product",
			Price:       15.99,
		},
		{
			SKU:         "LIST003",
			Name:        "Product 3",
			Description: "Third product",
			Price:       25.99,
		},
	}

	for _, p := range products {
		_, err := service.CreateProduct(ctx, p)
		if err != nil {
			t.Fatalf("Expected no error creating product, got %v", err)
		}
	}

	// List all products
	retrieved, err := service.ListProducts(ctx)
	if err != nil {
		t.Fatalf("Expected no error listing products, got %v", err)
	}

	if len(retrieved) != 3 {
		t.Fatalf("Expected 3 products, got %d", len(retrieved))
	}

	// Verify all products are present
	skus := make(map[string]bool)
	for _, p := range retrieved {
		skus[p.SKU] = true
	}

	for _, p := range products {
		if !skus[p.SKU] {
			t.Errorf("Product with SKU %s should be in the list", p.SKU)
		}
	}
}

func TestProductService_ListProducts_Empty(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()

	// List products when none exist
	retrieved, err := service.ListProducts(ctx)
	if err != nil {
		t.Fatalf("Expected no error listing empty products, got %v", err)
	}

	if len(retrieved) != 0 {
		t.Fatalf("Expected 0 products, got %d", len(retrieved))
	}
}

func TestProductService_CreateProduct_InvalidInput(t *testing.T) {
	repo := &MockProductRepository{
		products: make(map[string]*models.Product),
	}
	service := NewProductService(repo)

	ctx := context.Background()

	testCases := []struct {
		name    string
		req     *models.CreateProductRequest
		wantErr bool
	}{
		{
			name: "Empty SKU",
			req: &models.CreateProductRequest{
				SKU:         "",
				Name:        "Test Product",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: true,
		},
		{
			name: "Empty Name",
			req: &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "",
				Description: "A test product",
				Price:       9.99,
			},
			wantErr: true,
		},
		{
			name: "Negative Price",
			req: &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       -1.00,
			},
			wantErr: true,
		},
		{
			name: "Zero Price",
			req: &models.CreateProductRequest{
				SKU:         "TEST001",
				Name:        "Test Product",
				Description: "A test product",
				Price:       0.00,
			},
			wantErr: false, // Zero price should be allowed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := service.CreateProduct(ctx, tc.req)

			if tc.wantErr && err == nil {
				t.Fatalf("Expected error, got none")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
		})
	}
}
