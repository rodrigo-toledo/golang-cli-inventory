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

// StockMovementRepository provides methods for interacting with stock movement data in the database.
// It handles operations related to tracking stock movements such as additions, removals, and transfers.
type StockMovementRepository struct {
	queries *db.Queries
}

// NewStockMovementRepository creates a new instance of StockMovementRepository with the provided database queries.
func NewStockMovementRepository(queries *db.Queries) *StockMovementRepository {
	return &StockMovementRepository{
		queries: queries,
	}
}

func (r *StockMovementRepository) Create(ctx context.Context, movement *models.CreateStockMovementRequest) (*models.StockMovement, error) {
	params := db.CreateStockMovementParams{
		ProductID:      int32(movement.ProductID),
		FromLocationID: pgtype.Int4{Int32: int32(movement.SourceLocationID), Valid: movement.SourceLocationID != 0},
		ToLocationID:   pgtype.Int4{Int32: int32(movement.DestLocationID), Valid: movement.DestLocationID != 0},
		Quantity:       int32(movement.Quantity),
		MovementType:   movement.MovementType,
	}

	dbMovement, err := r.queries.CreateStockMovement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	var fromLocID, toLocID *int
	if dbMovement.FromLocationID.Valid {
		id := int(dbMovement.FromLocationID.Int32)
		fromLocID = &id
	}
	if dbMovement.ToLocationID.Valid {
		id := int(dbMovement.ToLocationID.Int32)
		toLocID = &id
	}

	return &models.StockMovement{
		ID:             int(dbMovement.ID),
		ProductID:      int(dbMovement.ProductID),
		FromLocationID: fromLocID,
		ToLocationID:   toLocID,
		Quantity:       int(dbMovement.Quantity),
		MovementType:   dbMovement.MovementType,
		CreatedAt:      dbMovement.CreatedAt.Time,
	}, nil
}

func (r *StockMovementRepository) ListByProductID(ctx context.Context, productID int) ([]models.StockMovement, error) {
	dbMovements, err := r.queries.ListStockMovementsByProduct(ctx, int32(productID))
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements by product id: %w", err)
	}

	movements := make([]models.StockMovement, len(dbMovements))
	for i, dbMovement := range dbMovements {
		var fromLocID, toLocID *int
		if dbMovement.FromLocationID.Valid {
			id := int(dbMovement.FromLocationID.Int32)
			fromLocID = &id
		}
		if dbMovement.ToLocationID.Valid {
			id := int(dbMovement.ToLocationID.Int32)
			toLocID = &id
		}
		movements[i] = models.StockMovement{
			ID:             int(dbMovement.ID),
			ProductID:      int(dbMovement.ProductID),
			FromLocationID: fromLocID,
			ToLocationID:   toLocID,
			Quantity:       int(dbMovement.Quantity),
			MovementType:   dbMovement.MovementType,
			CreatedAt:      dbMovement.CreatedAt.Time,
		}
	}

	return movements, nil
}
