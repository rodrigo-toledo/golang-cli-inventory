package service

import (
	"context"
	"fmt"

	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepositoryInterface interface {
	GetByID(ctx context.Context, id int) (*models.Location, error)
}

type StockRepositoryInterface interface {
	AddStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error)
	RemoveStock(ctx context.Context, productID, locationID, quantity int) (*models.Stock, error)
	GetLowStock(ctx context.Context, threshold int) ([]models.Stock, error)
}

type StockMovementRepositoryInterface interface {
	Create(ctx context.Context, movement *models.StockMovement) (*models.StockMovement, error)
}

type StockService struct {
	productRepo  ProductRepositoryInterface
	locationRepo LocationRepositoryInterface
	stockRepo    StockRepositoryInterface
	movementRepo StockMovementRepositoryInterface
	db           *pgxpool.Pool
}

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
