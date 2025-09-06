package service

import (
	"context"
	"fmt"

	"github.com/rodrigotoledo/cli-inventory/internal/models"
)

type ProductRepositoryInterface interface {
	Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	List(ctx context.Context) ([]models.Product, error)
}

type ProductService struct {
	repo ProductRepositoryInterface
}

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
