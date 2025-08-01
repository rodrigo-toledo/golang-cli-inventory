package repository

import (
	"context"
	"fmt"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
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
	// Handle nullable fields
	var fromLocationID, toLocationID pgtype.Int4
	if movement.FromLocationID != nil {
		fromLocationID = pgtype.Int4{Int32: int32(*movement.FromLocationID), Valid: true}
	}
	if movement.ToLocationID != nil {
		toLocationID = pgtype.Int4{Int32: int32(*movement.ToLocationID), Valid: true}
	}

	params := db.CreateStockMovementParams{
		ProductID:      int32(movement.ProductID),
		FromLocationID: fromLocationID,
		ToLocationID:   toLocationID,
		Quantity:       int32(movement.Quantity),
		MovementType:   movement.MovementType,
	}

	dbMovement, err := r.queries.CreateStockMovement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stock movement: %w", err)
	}

	// Convert pgtype.Int4 to *int
	var fromLoc, toLoc *int
	if dbMovement.FromLocationID.Valid {
		val := int(dbMovement.FromLocationID.Int32)
		fromLoc = &val
	}
	if dbMovement.ToLocationID.Valid {
		val := int(dbMovement.ToLocationID.Int32)
		toLoc = &val
	}

	return &models.StockMovement{
		ID:             int(dbMovement.ID),
		ProductID:      int(dbMovement.ProductID),
		FromLocationID: fromLoc,
		ToLocationID:   toLoc,
		Quantity:       int(dbMovement.Quantity),
		MovementType:   dbMovement.MovementType,
		CreatedAt:      dbMovement.CreatedAt.Time,
	}, nil
}

func (r *StockMovementRepository) List(ctx context.Context) ([]models.StockMovement, error) {
	dbMovements, err := r.queries.ListStockMovements(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock movements: %w", err)
	}

	movements := make([]models.StockMovement, len(dbMovements))
	for i, dbMovement := range dbMovements {
		// Convert pgtype.Int4 to *int
		var fromLoc, toLoc *int
		if dbMovement.FromLocationID.Valid {
			val := int(dbMovement.FromLocationID.Int32)
			fromLoc = &val
		}
		if dbMovement.ToLocationID.Valid {
			val := int(dbMovement.ToLocationID.Int32)
			toLoc = &val
		}

		movements[i] = models.StockMovement{
			ID:             int(dbMovement.ID),
			ProductID:      int(dbMovement.ProductID),
			FromLocationID: fromLoc,
			ToLocationID:   toLoc,
			Quantity:       int(dbMovement.Quantity),
			MovementType:   dbMovement.MovementType,
			CreatedAt:      dbMovement.CreatedAt.Time,
		}
	}

	return movements, nil
}
