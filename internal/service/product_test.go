package service

import (
	"context"
	"testing"

	"cli-inventory/internal/models"
)

// MockProductRepository is a mock implementation of ProductRepositoryInterface for testing
type MockProductRepository struct {
	products map[string]*models.Product
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
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
