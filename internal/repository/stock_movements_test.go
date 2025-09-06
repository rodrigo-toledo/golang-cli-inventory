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

		fromLocationID := 1
		toLocationID := 2
		movement := &models.StockMovement{
			ProductID:      1,
			FromLocationID: &fromLocationID,
			ToLocationID:   &toLocationID,
			Quantity:       10,
			MovementType:   "MOVE",
		}

		expectedMovement := db.StockMovement{
			ID:             1,
			ProductID:      1,
			FromLocationID: pgtype.Int4{Int32: 1, Valid: true},
			ToLocationID:   pgtype.Int4{Int32: 2, Valid: true},
			Quantity:       10,
			MovementType:   "MOVE",
			CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		// Mock the QueryRow method
		mockRow := new(MockRow) // This will use the MockRow from locations_test.go
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*int32)
				*arg = expectedMovement.ID
				arg1 := args.Get(1).(*int32)
				*arg1 = expectedMovement.ProductID
				arg2 := args.Get(2).(*pgtype.Int4)
				*arg2 = expectedMovement.FromLocationID
				arg3 := args.Get(3).(*pgtype.Int4)
				*arg3 = expectedMovement.ToLocationID
				arg4 := args.Get(4).(*int32)
				*arg4 = expectedMovement.Quantity
				arg5 := args.Get(5).(*string)
				*arg5 = expectedMovement.MovementType
				arg6 := args.Get(6).(*pgtype.Timestamptz)
				*arg6 = expectedMovement.CreatedAt
			})

		mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

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
		mockRow.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		fromLocationID := 1
		toLocationID := 2
		movement := &models.StockMovement{
			ProductID:      1,
			FromLocationID: &fromLocationID,
			ToLocationID:   &toLocationID,
			Quantity:       10,
			MovementType:   "MOVE",
		}

		// Mock the QueryRow method to return an error
		mockRow := new(MockRow) // This will use the MockRow from locations_test.go
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("database error"))

		mockDB.On("QueryRow", mock.Anything, mock.AnythingOfType("string"), mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mockRow)

		result, err := repo.Create(context.Background(), movement)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to create stock movement: database error")

		mockDB.AssertExpectations(t)
		mockRow.AssertExpectations(t)
	})
}

func TestStockMovementRepository_List(t *testing.T) {
	expectedMovements := []db.StockMovement{
		{
			ID:             1,
			ProductID:      1,
			FromLocationID: pgtype.Int4{Int32: 1, Valid: true},
			ToLocationID:   pgtype.Int4{Int32: 2, Valid: true},
			Quantity:       10,
			MovementType:   "MOVE",
			CreatedAt:      pgtype.Timestamptz{Time: time.Now(), Valid: true},
		},
	}

	t.Run("successful list", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		mockRows := new(MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*int32)
			*arg = expectedMovements[0].ID
			arg1 := args.Get(1).(*int32)
			*arg1 = expectedMovements[0].ProductID
			arg2 := args.Get(2).(*pgtype.Int4)
			*arg2 = expectedMovements[0].FromLocationID
			arg3 := args.Get(3).(*pgtype.Int4)
			*arg3 = expectedMovements[0].ToLocationID
			arg4 := args.Get(4).(*int32)
			*arg4 = expectedMovements[0].Quantity
			arg5 := args.Get(5).(*string)
			*arg5 = expectedMovements[0].MovementType
			arg6 := args.Get(6).(*pgtype.Timestamptz)
			*arg6 = expectedMovements[0].CreatedAt
		}).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Once()

		mockDB.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, nil)

		result, err := repo.List(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedMovements[0].ID, int32(result[0].ID))

		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		mockDB := new(MockDBTXForStock)
		queries := db.New(mockDB)
		repo := NewStockMovementRepository(queries)

		mockRows := new(MockRows)
		mockDB.On("Query", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(mockRows, errors.New("database error"))

		result, err := repo.List(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, "failed to list stock movements: database error")

		mockDB.AssertExpectations(t)
	})
}
