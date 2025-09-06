package repository

import (
	"context"
	"fmt"

	"github.com/rodrigotoledo/cli-inventory/internal/db"
	"github.com/rodrigotoledo/cli-inventory/internal/models"
)

type ProductRepository struct {
	queries *db.Queries
}

func NewProductRepository(queries *db.Queries) *ProductRepository {
	return &ProductRepository{
		queries: queries,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
	params := db.CreateProductParams{
		Sku:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}

	dbProduct, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
		CreatedAt:   dbProduct.CreatedAt,
	}, nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
		CreatedAt:   dbProduct.CreatedAt,
	}, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Price:       dbProduct.Price,
		CreatedAt:   dbProduct.CreatedAt,
	}, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	dbProducts, err := r.queries.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]models.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = models.Product{
			ID:          int(dbProduct.ID),
			SKU:         dbProduct.Sku,
			Name:        dbProduct.Name,
			Description: dbProduct.Description,
			Price:       dbProduct.Price,
			CreatedAt:   dbProduct.CreatedAt,
		}
	}

	return products, nil
}
