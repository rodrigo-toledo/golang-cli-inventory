package repository

import (
	"context"
	"fmt"

	"github.com/rodrigotoledo/cli-inventory/internal/db"
	"github.com/rodrigotoledo/cli-inventory/internal/models"
)

type StockMovementRepository struct {
	queries *db.Queries
}

func NewStockMovementRepository(queries *db.Queries) *StockMovementRepository {
	return &StockMovementRepository{
		queries: queries,
	}
}

func (r *StockMovementRepository) Create(ctx context.Context, movement *models.StockMovement) (*models.StockMovement, error) {
	params := db.CreateStockMovementParams{
		ProductID:    int32(movement.ProductID),
		Quantity:     int32(movement.Quantity),
		MovementType: movement.MovementType,
	}

	// Handle nullable fields
	if movement.FromLocationID != nil {
		fromID := int32(*movement.FromLocationID)
		params.FromLocationID = &fromID
	}
	if movement.ToLocationID != nil {
		toID := int32(*movement.ToLocationID)
		params.ToLocationID = &toID
	}

	dbMovement, err := r.queries.CreateStockMovement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	return &models.StockMovement{
		ID:             int(dbMovement.ID),
		ProductID:      int(dbMovement.ProductID),
		FromLocationID: (*int)(params.FromLocationID),
		ToLocationID:   (*int)(params.ToLocationID),
		Quantity:       int(dbMovement.Quantity),
		MovementType:   dbMovement.MovementType,
		CreatedAt:      dbMovement.CreatedAt,
	}, nil
}

func (r *StockMovementRepository) List(ctx context.Context) ([]models.StockMovement, error) {
	dbMovements, err := r.queries.ListStockMovements(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements: %w", err)
	}

	movements := make([]models.StockMovement, len(dbMovements))
	for i, dbMovement := range dbMovements {
		movements[i] = models.StockMovement{
			ID:             int(dbMovement.ID),
			ProductID:      int(dbMovement.ProductID),
			FromLocationID: nil, // These might need to be handled differently
			ToLocationID:   nil, // These might need to be handled differently
			Quantity:       int(dbMovement.Quantity),
			MovementType:   dbMovement.MovementType,
			CreatedAt:      dbMovement.CreatedAt,
		}
	}

	return movements, nil
}
