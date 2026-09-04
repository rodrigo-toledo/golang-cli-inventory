package repository

import (
	"context"
	"errors"
	"strings"
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

		mockRow := new(MockRowForStock)
		mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
			return strings.Contains(query, "INSERT INTO stock_movements")
		}), mock.AnythingOfType("[]interface {}")).Return(mockRow)

		mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
			*(args.Get(0).(*int32)) = expectedMovement.ID
			*(args.Get(1).(*int32)) = expectedMovement.ProductID
			*(args.Get(2).(*pgtype.Int4)) = expectedMovement.FromLocationID
			*(args.Get(3).(*pgtype.Int4)) = expectedMovement.ToLocationID
			*(args.Get(4).(*int32)) = expectedMovement.Quantity
			*(args.Get(5).(*string)) = expectedMovement.MovementType
			*(args.Get(6).(*pgtype.Timestamptz)) = expectedMovement.CreatedAt
		})

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

		mockRow := new(MockRowForStock)
		mockDB.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
			return strings.Contains(query, "INSERT INTO stock_movements")
		}), mock.AnythingOfType("[]interface {}")).Return(mockRow)

		mockRow.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(errors.New("database error"))

		result, err := repo.Create(context.Background(), movement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to create stock movement: database error")

		mockDB.AssertExpectations(t)
		mockRow.AssertExpectations(t)
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

		mockRows := new(MockRowsForStock)
		mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int32"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*pgtype.Int4"), mock.AnythingOfType("*int32"), mock.AnythingOfType("*string"), mock.AnythingOfType("*pgtype.Timestamptz")).Return(nil).Run(func(args mock.Arguments) {
			*(args.Get(0).(*int32)) = expectedMovements[0].ID
			*(args.Get(1).(*int32)) = expectedMovements[0].ProductID
			*(args.Get(2).(*pgtype.Int4)) = expectedMovements[0].FromLocationID
			*(args.Get(3).(*pgtype.Int4)) = expectedMovements[0].ToLocationID
			*(args.Get(4).(*int32)) = expectedMovements[0].Quantity
			*(args.Get(5).(*string)) = expectedMovements[0].MovementType
			*(args.Get(6).(*pgtype.Timestamptz)) = expectedMovements[0].CreatedAt
		})
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()
		mockRows.On("Err").Return(nil)

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

		mockDB.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("database error"))

		result, err := repo.ListByProductID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to list stock movements by product id: database error")

		mockDB.AssertExpectations(t)
	})
}
