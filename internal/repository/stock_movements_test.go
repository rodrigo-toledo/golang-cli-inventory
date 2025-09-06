package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStockMovementRepository_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		movement := &models.CreateStockMovementRequest{
			ProductID:        1,
			SourceLocationID: 1,
			DestLocationID:   2,
			Quantity:         10,
		}

		expectedMovement := db.StockMovement{
			ID:             1,
			ProductID:      1,
			FromLocationID: pgtype.Int4{Int32: 1, Valid: true},
			ToLocationID:   pgtype.Int4{Int32: 2, Valid: true},
			Quantity:       10,
			MovementType:   "transfer",
			CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		mockDB.On("CreateStockMovement", mock.Anything, mock.Anything).Return(expectedMovement, nil)

		result, err := repo.Create(context.Background(), movement)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedMovement.ID, int32(result.ID))
		assert.Equal(t, expectedMovement.ProductID, int32(result.ProductID))
		assert.Equal(t, expectedMovement.FromLocationID.Int32, int32(*result.FromLocationID))
		assert.Equal(t, expectedMovement.ToLocationID.Int32, int32(*result.ToLocationID))
		assert.Equal(t, expectedMovement.Quantity, int32(result.Quantity))
		assert.Equal(t, expectedMovement.MovementType, result.MovementType)

		mockDB.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		movement := &models.CreateStockMovementRequest{
			ProductID:        1,
			SourceLocationID: 1,
			DestLocationID:   2,
			Quantity:         10,
		}

		mockDB.On("CreateStockMovement", mock.Anything, mock.Anything).Return(db.StockMovement{}, errors.New("database error"))

		result, err := repo.Create(context.Background(), movement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to create stock movement: database error")

		mockDB.AssertExpectations(t)
	})
}

func TestStockMovementRepository_ListByProductID(t *testing.T) {
	expectedMovements := []db.StockMovement{
		{
			ID:             1,
			ProductID:      1,
			FromLocationID: pgtype.Int4{Int32: 1, Valid: true},
			ToLocationID:   pgtype.Int4{Int32: 2, Valid: true},
			Quantity:       10,
			MovementType:   "transfer",
			CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	t.Run("successful list", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		mockDB.On("ListStockMovementsByProduct", mock.Anything, mock.Anything).Return(expectedMovements, nil)

		result, err := repo.ListByProductID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedMovements[0].ID, int32(result[0].ID))

		mockDB.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		mockDB.On("ListStockMovementsByProduct", mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

		result, err := repo.ListByProductID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to list stock movements by product id: database error")

		mockDB.AssertExpectations(t)
	})
}
