// Package service provides business logic implementations for the inventory management system.
// It contains services that handle the core functionality such as product management,
// stock management, and location management.
package service

import (
	"context"
	"errors"
	"fmt"

	"cli-inventory/internal/models"
)

// ErrProductNotFound is returned when a product cannot be found by its SKU or ID.
var ErrProductNotFound = errors.New("product not found")

// ProductService provides methods for managing products in the inventory system.
// It handles operations such as creating products, retrieving product information,
// and listing all products.
type ProductService struct {
	repo ProductRepositoryInterface
}

// NewProductService creates a new instance of ProductService with the provided product repository.
func NewProductService(repo ProductRepositoryInterface) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error) {
	// Check if product with this SKU already exists
	existing, err := s.repo.GetBySKU(ctx, req.SKU)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("product with SKU %s already exists", req.SKU)
	}

	// Create the product
	product, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) GetProductBySKU(ctx context.Context, sku string) (*models.Product, error) {
	product, err := s.repo.GetBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *ProductService) ListProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}