// Package service provides business logic implementations for the inventory management system.
// It contains services that handle the core functionality such as product management,
// stock management, and location management.
package service

import (
	"context"

	"cli-inventory/internal/models"
)

// ProductRepositoryInterface defines the contract for product data access operations.
// It specifies the methods that any product repository implementation must provide.
type ProductRepositoryInterface interface {
	Create(ctx context.Context, product *models.CreateProductRequest) (*models.Product, error)
	GetBySKU(ctx context.Context, sku string) (*models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	List(ctx context.Context) ([]models.Product, error)
}

// LocationRepositoryInterface defines the contract for location data access operations.
// It specifies the methods that any location repository implementation must provide.
type LocationRepositoryInterface interface {
	Create(ctx context.Context, location *models.CreateLocationRequest) (*models.Location, error)
	GetByName(ctx context.Context, name string) (*models.Location, error)
	GetByID(ctx context.Context, id int) (*models.Location, error)
	List(ctx context.Context) ([]models.Location, error)
}

// StockRepositoryInterface defines the contract for stock data access operations.
// It specifies the methods that any stock repository implementation must provide.
type StockRepositoryInterface interface {
	AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error)
	RemoveStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error)
	GetLowStock(ctx context.Context, threshold int) ([]models.Stock, error)
	GetByProductAndLocation(ctx context.Context, productID, locationID int) (*models.Stock, error)
}

// StockMovementRepositoryInterface defines the contract for stock movement data access operations.
// It specifies the methods that any stock movement repository implementation must provide.
type StockMovementRepositoryInterface interface {
	Create(ctx context.Context, movement *models.StockMovement) (*models.StockMovement, error)
}

// ProductServiceInterface defines the contract for product business logic operations.
// It specifies the methods that any product service implementation must provide.
type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, req *models.CreateProductRequest) (*models.Product, error)
	GetProductBySKU(ctx context.Context, sku string) (*models.Product, error)
	ListProducts(ctx context.Context) ([]models.Product, error)
}

// LocationServiceInterface defines the contract for location business logic operations.
// It specifies the methods that any location service implementation must provide.
type LocationServiceInterface interface {
	CreateLocation(ctx context.Context, req *models.CreateLocationRequest) (*models.Location, error)
	GetLocationByName(ctx context.Context, name string) (*models.Location, error)
	ListLocations(ctx context.Context) ([]models.Location, error)
}

// StockServiceInterface defines the contract for stock business logic operations.
// It specifies the methods that any stock service implementation must provide.
type StockServiceInterface interface {
	AddStock(ctx context.Context, req *models.AddStockRequest) (*models.Stock, error)
	MoveStock(ctx context.Context, req *models.MoveStockRequest) (*models.Stock, error)
	GetLowStockReport(ctx context.Context, threshold int) ([]models.Stock, error)
}