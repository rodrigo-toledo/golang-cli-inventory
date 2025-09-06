// Package service provides business logic implementations for the inventory management system.
// It contains services that handle the core functionality such as product management,
// stock management, and location management.
package service

import (
	"context"
	"fmt"

	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LocationRepositoryInterface defines the contract for location data access operations.
// It specifies the methods that any location repository implementation must provide.
type LocationRepositoryInterface interface {
	GetByID(ctx context.Context, id int) (*models.Location, error)
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

// StockService provides methods for managing stock levels and movements in the inventory system.
// It handles operations such as adding stock, moving stock between locations, and generating reports.
type StockService struct {
	productRepo  ProductRepositoryInterface
	locationRepo LocationRepositoryInterface
	stockRepo    StockRepositoryInterface
	movementRepo StockMovementRepositoryInterface
	db           *pgxpool.Pool
}

// NewStockService creates a new instance of StockService with the provided repositories and database connection.
func NewStockService(
	productRepo ProductRepositoryInterface,
	locationRepo LocationRepositoryInterface,
	stockRepo StockRepositoryInterface,
	movementRepo StockMovementRepositoryInterface,
	db *pgxpool.Pool,
) *StockService {
	return &StockService{
		productRepo:  productRepo,
		locationRepo: locationRepo,
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
		db:           db,
	}
}

func (s *StockService) AddStock(ctx context.Context, req *models.AddStockRequest) (*models.Stock, error) {
	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %d does not exist", req.ProductID)
	}

	// Check if location exists
	_, err = s.locationRepo.GetByID(ctx, req.LocationID)
	if err != nil {
		return nil, fmt.Errorf("location with ID %d does not exist", req.LocationID)
	}

	// Add stock
	stock, err := s.stockRepo.AddStock(ctx, req.ProductID, req.LocationID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	// Record the movement
	movement := &models.StockMovement{
		ProductID:    req.ProductID,
		ToLocationID: &req.LocationID,
		Quantity:     req.Quantity,
		MovementType: "ADD",
	}
	_, err = s.movementRepo.Create(ctx, movement)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to record stock movement: %v\n", err)
	}

	return stock, nil
}

func (s *StockService) MoveStock(ctx context.Context, req *models.MoveStockRequest) (*models.Stock, error) {
	// Validate input
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	if req.FromLocationID == req.ToLocationID {
		return nil, fmt.Errorf("source and destination locations cannot be the same")
	}

	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product with ID %d does not exist", req.ProductID)
	}

	// Check if from location exists
	_, err = s.locationRepo.GetByID(ctx, req.FromLocationID)
	if err != nil {
		return nil, fmt.Errorf("from location with ID %d does not exist", req.FromLocationID)
	}

	// Check if to location exists
	_, err = s.locationRepo.GetByID(ctx, req.ToLocationID)
	if err != nil {
		return nil, fmt.Errorf("to location with ID %d does not exist", req.ToLocationID)
	}

	// Check if there's sufficient stock at the source location
	currentStock, err := s.stockRepo.GetByProductAndLocation(ctx, req.ProductID, req.FromLocationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check current stock: %w", err)
	}

	if currentStock.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: only %d available, requested %d", currentStock.Quantity, req.Quantity)
	}

	// If db is nil (e.g., in tests), perform operations without transaction
	if s.db == nil {
		// Remove stock from source location
		_, err = s.stockRepo.RemoveStock(ctx, req.ProductID, req.FromLocationID, req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to remove stock from source location: %w", err)
		}

		// Add stock to destination location
		stock, err := s.stockRepo.AddStock(ctx, req.ProductID, req.ToLocationID, req.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to add stock to destination location: %w", err)
		}

		// Record the movement
		movement := &models.StockMovement{
			ProductID:      req.ProductID,
			FromLocationID: &req.FromLocationID,
			ToLocationID:   &req.ToLocationID,
			Quantity:       req.Quantity,
			MovementType:   "MOVE",
		}
		_, err = s.movementRepo.Create(ctx, movement)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Warning: failed to record stock movement: %v\n", err)
		}

		return stock, nil
	}

	// Start transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Remove stock from source location
	_, err = s.stockRepo.RemoveStock(ctx, req.ProductID, req.FromLocationID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to remove stock from source location: %w", err)
	}

	// Add stock to destination location
	stock, err := s.stockRepo.AddStock(ctx, req.ProductID, req.ToLocationID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add stock to destination location: %w", err)
	}

	// Record the movement
	movement := &models.StockMovement{
		ProductID:      req.ProductID,
		FromLocationID: &req.FromLocationID,
		ToLocationID:   &req.ToLocationID,
		Quantity:       req.Quantity,
		MovementType:   "MOVE",
	}
	_, err = s.movementRepo.Create(ctx, movement)
	if err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to record stock movement: %v\n", err)
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return stock, nil
}

func (s *StockService) GetLowStockReport(ctx context.Context, threshold int) ([]models.Stock, error) {
	stocks, err := s.stockRepo.GetLowStock(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get low stock report: %w", err)
	}
	return stocks, nil
}
