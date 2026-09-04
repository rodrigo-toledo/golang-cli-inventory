// Package service provides business logic implementations for the inventory management system.
// It contains services that handle the core functionality such as product management,
// stock management, and location management.
package service

import (
	"context"
	"errors"
	"fmt"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"
	"cli-inventory/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Movement types recorded in the stock_movements table.
const (
	movementTypeAdd  = "ADD"
	movementTypeMove = "MOVE"
)

// ErrInsufficientStock is returned when an attempt is made to move more stock than is available.
var ErrInsufficientStock = errors.New("insufficient stock")

// StockService provides methods for managing stock levels and movements in the inventory system.
// It handles operations such as adding stock, moving stock between locations, and generating reports.
type StockService struct {
	productRepo  ProductRepositoryInterface
	locationRepo LocationRepositoryInterface
	stockRepo    StockRepositoryInterface
	movementRepo StockMovementRepositoryInterface
	pool         *pgxpool.Pool
}

// NewStockService creates a new instance of StockService with the provided repositories and database connection.
func NewStockService(
	productRepo ProductRepositoryInterface,
	locationRepo LocationRepositoryInterface,
	stockRepo StockRepositoryInterface,
	movementRepo StockMovementRepositoryInterface,
	pool *pgxpool.Pool,
) *StockService {
	return &StockService{
		productRepo:  productRepo,
		locationRepo: locationRepo,
		stockRepo:    stockRepo,
		movementRepo: movementRepo,
		pool:         pool,
	}
}

// AddStock increases the stock level of a product at a location and records an ADD movement.
// When a database pool is available the stock update and the movement record are committed
// atomically; otherwise the provided repositories are used directly (unit-test mode).
func (s *StockService) AddStock(ctx context.Context, req *models.AddStockRequest) (*models.Stock, error) {
	if _, err := s.productRepo.GetByID(ctx, req.ProductID); err != nil {
		return nil, fmt.Errorf("product with ID %d does not exist", req.ProductID)
	}

	if _, err := s.locationRepo.GetByID(ctx, req.LocationID); err != nil {
		return nil, fmt.Errorf("location with ID %d does not exist", req.LocationID)
	}

	if s.pool == nil {
		return s.addStockWithRepos(ctx, s.stockRepo, s.movementRepo, req)
	}

	var stock *models.Stock
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		q := db.New(tx)
		var err error
		stock, err = s.addStockWithRepos(ctx, repository.NewStockRepository(q), repository.NewStockMovementRepository(q), req)
		return err
	})
	if err != nil {
		return nil, err
	}
	return stock, nil
}

// addStockWithRepos performs the stock increase and records the ADD movement using the
// given repositories. Any error aborts the surrounding transaction, if one is active.
func (s *StockService) addStockWithRepos(
	ctx context.Context,
	stockRepo StockRepositoryInterface,
	movementRepo StockMovementRepositoryInterface,
	req *models.AddStockRequest,
) (*models.Stock, error) {
	stock, err := stockRepo.AddStock(ctx, req.ProductID, req.LocationID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	movement := &models.CreateStockMovementRequest{
		ProductID:      req.ProductID,
		DestLocationID: req.LocationID,
		Quantity:       req.Quantity,
		MovementType:   movementTypeAdd,
	}
	if _, err := movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to record stock movement: %w", err)
	}

	return stock, nil
}

// MoveStock transfers stock between locations and records a MOVE movement. The stock
// decrement, increment, and movement record are committed atomically when a database
// pool is available.
func (s *StockService) MoveStock(ctx context.Context, req *models.MoveStockRequest) (*models.Stock, error) {
	if req.Quantity <= 0 {
		return nil, fmt.Errorf("quantity must be positive")
	}

	if req.FromLocationID == req.ToLocationID {
		return nil, fmt.Errorf("source and destination locations cannot be the same")
	}

	if _, err := s.productRepo.GetByID(ctx, req.ProductID); err != nil {
		return nil, fmt.Errorf("product with ID %d does not exist", req.ProductID)
	}

	if _, err := s.locationRepo.GetByID(ctx, req.FromLocationID); err != nil {
		return nil, fmt.Errorf("from location with ID %d does not exist", req.FromLocationID)
	}

	if _, err := s.locationRepo.GetByID(ctx, req.ToLocationID); err != nil {
		return nil, fmt.Errorf("to location with ID %d does not exist", req.ToLocationID)
	}

	if s.pool == nil {
		return s.moveStockWithRepos(ctx, s.stockRepo, s.movementRepo, req)
	}

	var stock *models.Stock
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		q := db.New(tx)
		var err error
		stock, err = s.moveStockWithRepos(ctx, repository.NewStockRepository(q), repository.NewStockMovementRepository(q), req)
		return err
	})
	if err != nil {
		return nil, err
	}
	return stock, nil
}

// moveStockWithRepos performs the stock transfer and records the MOVE movement using the
// given repositories. The source decrement is guarded by the SQL predicate
// quantity >= requested, so pgx.ErrNoRows signals insufficient stock. Any error aborts
// the surrounding transaction, if one is active.
func (s *StockService) moveStockWithRepos(
	ctx context.Context,
	stockRepo StockRepositoryInterface,
	movementRepo StockMovementRepositoryInterface,
	req *models.MoveStockRequest,
) (*models.Stock, error) {
	if _, err := stockRepo.RemoveStock(ctx, req.ProductID, req.FromLocationID, req.Quantity); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: product %d has fewer than %d at location %d",
				ErrInsufficientStock, req.ProductID, req.Quantity, req.FromLocationID)
		}
		return nil, fmt.Errorf("failed to remove stock from source location: %w", err)
	}

	stock, err := stockRepo.AddStock(ctx, req.ProductID, req.ToLocationID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add stock to destination location: %w", err)
	}

	movement := &models.CreateStockMovementRequest{
		ProductID:        req.ProductID,
		SourceLocationID: req.FromLocationID,
		DestLocationID:   req.ToLocationID,
		Quantity:         req.Quantity,
		MovementType:     movementTypeMove,
	}
	if _, err := movementRepo.Create(ctx, movement); err != nil {
		return nil, fmt.Errorf("failed to record stock movement: %w", err)
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
