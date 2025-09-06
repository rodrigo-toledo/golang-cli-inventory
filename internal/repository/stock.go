// Package repository provides data access implementations for the inventory management system.
// It contains repositories that handle interactions with the database for products, locations,
// stock, and stock movements.
package repository

import (
	"context"
	"fmt"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"
)

// StockRepository provides methods for interacting with stock data in the database.
// It handles operations related to stock levels and stock movements.
type StockRepository struct {
	queries *db.Queries
}

// NewStockRepository creates a new instance of StockRepository with the provided database queries.
func NewStockRepository(queries *db.Queries) *StockRepository {
	return &StockRepository{
		queries: queries,
	}
}

func (r *StockRepository) Create(ctx context.Context, stock *models.AddStockRequest) (*models.Stock, error) {
	params := db.CreateStockParams{
		ProductID:  int32(stock.ProductID),
		LocationID: int32(stock.LocationID),
		Quantity:   int32(stock.Quantity),
	}

	dbStock, err := r.queries.CreateStock(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock: %w", err)
	}

	return &models.Stock{
		ID:         int(dbStock.ID),
		ProductID:  int(dbStock.ProductID),
		LocationID: int(dbStock.LocationID),
		Quantity:   int(dbStock.Quantity),
		CreatedAt:  dbStock.CreatedAt.Time,
		UpdatedAt:  dbStock.UpdatedAt.Time,
	}, nil
}

func (r *StockRepository) GetByProductAndLocation(ctx context.Context, productID, locationID int) (*models.Stock, error) {
	params := db.GetStockByProductAndLocationParams{
		ProductID:  int32(productID),
		LocationID: int32(locationID),
	}

	dbStock, err := r.queries.GetStockByProductAndLocation(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}

	return &models.Stock{
		ID:         int(dbStock.ID),
		ProductID:  int(dbStock.ProductID),
		LocationID: int(dbStock.LocationID),
		Quantity:   int(dbStock.Quantity),
		CreatedAt:  dbStock.CreatedAt.Time,
		UpdatedAt:  dbStock.UpdatedAt.Time,
	}, nil
}

func (r *StockRepository) AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	params := db.AddStockParams{
		ProductID:  int32(productID),
		LocationID: int32(locationID),
		Quantity:   int32(quantity),
	}

	dbStock, err := r.queries.AddStock(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	return &models.Stock{
		ID:         int(dbStock.ID),
		ProductID:  int(dbStock.ProductID),
		LocationID: int(dbStock.LocationID),
		Quantity:   int(dbStock.Quantity),
		CreatedAt:  dbStock.CreatedAt.Time,
		UpdatedAt:  dbStock.UpdatedAt.Time,
	}, nil
}

func (r *StockRepository) RemoveStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error) {
	params := db.RemoveStockParams{
		ProductID:  int32(productID),
		LocationID: int32(locationID),
		Quantity:   int32(quantity),
	}

	dbStock, err := r.queries.RemoveStock(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}

	return &models.Stock{
		ID:         int(dbStock.ID),
		ProductID:  int(dbStock.ProductID),
		LocationID: int(dbStock.LocationID),
		Quantity:   int(dbStock.Quantity),
		CreatedAt:  dbStock.CreatedAt.Time,
		UpdatedAt:  dbStock.UpdatedAt.Time,
	}, nil
}

func (r *StockRepository) GetLowStock(ctx context.Context, threshold int) ([]models.Stock, error) {
	dbStocks, err := r.queries.GetStockByProduct(ctx, int32(threshold))
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock: %w", err)
	}

	stocks := make([]models.Stock, len(dbStocks))
	for i, dbStock := range dbStocks {
		stocks[i] = models.Stock{
			ID:         int(dbStock.ID),
			ProductID:  int(dbStock.ProductID),
			LocationID: int(dbStock.LocationID),
			Quantity:   int(dbStock.Quantity),
			CreatedAt:  dbStock.CreatedAt.Time,
			UpdatedAt:  dbStock.UpdatedAt.Time,
		}
	}

	return stocks, nil
}
