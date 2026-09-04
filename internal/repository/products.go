// Package repository provides data access implementations for the inventory management system.
// It contains repositories that handle interactions with the database for products, locations,
// stock, and stock movements.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5"
	pgtype "github.com/jackc/pgx/v5/pgtype"
)

// ProductRepository provides methods for interacting with product data in the database.
// It implements the ProductRepositoryInterface defined in the service package.
type ProductRepository struct {
	queries *db.Queries
}

// NewProductRepository creates a new instance of ProductRepository with the provided database queries.
func NewProductRepository(queries *db.Queries) *ProductRepository {
	return &ProductRepository{
		queries: queries,
	}
}

func (r *ProductRepository) Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error) {
	// Convert string to pgtype.Text
	description := pgtype.Text{String: product.Description, Valid: true}

	// Handle price conversion
	price := pgtype.Numeric{}
	if product.Price >= 0 {
		if err := price.Scan(strconv.FormatFloat(product.Price, 'f', -1, 64)); err != nil {
			return nil, fmt.Errorf("invalid price %v: %w", product.Price, err)
		}
		price.Valid = true
	}

	params := db.CreateProductParams{
		Sku:         product.SKU,
		Name:        product.Name,
		Description: description,
		Price:       price,
	}

	dbProduct, err := r.queries.CreateProduct(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return mapDBProductToModel(&dbProduct), nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductBySKU(ctx, sku)
	if err != nil {
		// If no product is found, return nil instead of an error
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	return mapDBProductToModel(&dbProduct), nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, int32(id))
	if err != nil {
		// If no product is found, return nil instead of an error
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	return mapDBProductToModel(&dbProduct), nil
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	dbProducts, err := r.queries.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := mapDBProductsToModels(dbProducts)

	return products, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteProduct(ctx, int32(id))
}
