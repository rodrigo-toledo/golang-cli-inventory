// Package repository provides data access implementations for the inventory management system.
// It contains repositories that handle interactions with the database for products, locations,
// stock, and stock movements.
package repository

import (
	"context"
	"fmt"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
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
	description := pgtype.Text{}
	description.Scan(product.Description)

	// Convert float64 to pgtype.Numeric
	price := pgtype.Numeric{}
	price.Scan(product.Price)

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

	// Convert pgtype.Text to string
	descriptionStr := ""
	if dbProduct.Description.Valid {
		descriptionStr = dbProduct.Description.String
	}

	// Convert pgtype.Numeric to float64
	var priceFloat float64
	if dbProduct.Price.Valid {
		if val, err := dbProduct.Price.Value(); err == nil {
			if floatVal, ok := val.(float64); ok {
				priceFloat = floatVal
			}
		}
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: descriptionStr,
		Price:       priceFloat,
		CreatedAt:   dbProduct.CreatedAt.Time,
	}, nil
}

func (r *ProductRepository) GetBySKU(ctx context.Context, sku string) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %w", err)
	}

	// Convert pgtype.Text to string
	descriptionStr := ""
	if dbProduct.Description.Valid {
		descriptionStr = dbProduct.Description.String
	}

	// Convert pgtype.Numeric to float64
	var priceFloat float64
	if dbProduct.Price.Valid {
		if val, err := dbProduct.Price.Value(); err == nil {
			if floatVal, ok := val.(float64); ok {
				priceFloat = floatVal
			}
		}
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: descriptionStr,
		Price:       priceFloat,
		CreatedAt:   dbProduct.CreatedAt.Time,
	}, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	dbProduct, err := r.queries.GetProductByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get product by ID: %w", err)
	}

	// Convert pgtype.Text to string
	descriptionStr := ""
	if dbProduct.Description.Valid {
		descriptionStr = dbProduct.Description.String
	}

	// Convert pgtype.Numeric to float64
	var priceFloat float64
	if dbProduct.Price.Valid {
		if val, err := dbProduct.Price.Value(); err == nil {
			if floatVal, ok := val.(float64); ok {
				priceFloat = floatVal
			}
		}
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: descriptionStr,
		Price:       priceFloat,
		CreatedAt:   dbProduct.CreatedAt.Time,
	}, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	dbProducts, err := r.queries.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	products := make([]models.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		// Convert pgtype.Text to string
		descriptionStr := ""
		if dbProduct.Description.Valid {
			descriptionStr = dbProduct.Description.String
		}

		// Convert pgtype.Numeric to float64
		var priceFloat float64
		if dbProduct.Price.Valid {
			if val, err := dbProduct.Price.Value(); err == nil {
				if floatVal, ok := val.(float64); ok {
					priceFloat = floatVal
				}
			}
		}

		products[i] = models.Product{
			ID:          int(dbProduct.ID),
			SKU:         dbProduct.Sku,
			Name:        dbProduct.Name,
			Description: descriptionStr,
			Price:       priceFloat,
			CreatedAt:   dbProduct.CreatedAt.Time,
		}
	}

	return products, nil
}
