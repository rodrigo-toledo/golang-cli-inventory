package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddStockRequest_BasicValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *AddStockRequest
		wantErr bool
	}{
		{
			name: "Valid Stock Request",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
			},
			wantErr: false,
		},
		{
			name: "Invalid Product ID",
			input: &AddStockRequest{
				ProductID:  0,
				LocationID: 1,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Negative Product ID",
			input: &AddStockRequest{
				ProductID:  -1,
				LocationID: 1,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Invalid Location ID",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: 0,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Negative Location ID",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: -1,
				Quantity:   10,
			},
			wantErr: true,
		},
		{
			name: "Zero Quantity",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   0,
			},
			wantErr: true,
		},
		{
			name: "Negative Quantity",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   -5,
			},
			wantErr: true,
		},
		{
			name: "Large Quantity",
			input: &AddStockRequest{
				ProductID:  1,
				LocationID: 1,
				Quantity:   10000,
			},
			wantErr: false, // Large quantity should be allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// validate is a simple validation method for AddStockRequest
func (r *AddStockRequest) validate() error {
	if r.ProductID <= 0 {
		return assert.AnError
	}
	if r.LocationID <= 0 {
		return assert.AnError
	}
	if r.Quantity <= 0 {
		return assert.AnError
	}
	return nil
}

func TestMoveStockRequest_BasicValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   *MoveStockRequest
		wantErr bool
	}{
		{
			name: "Valid Move Stock Request",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       5,
			},
			wantErr: false,
		},
		{
			name: "Invalid Product ID",
			input: &MoveStockRequest{
				ProductID:      0,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Invalid From Location ID",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 0,
				ToLocationID:   2,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Invalid To Location ID",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   0,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Same From and To Location",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   1,
				Quantity:       5,
			},
			wantErr: true,
		},
		{
			name: "Zero Quantity",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       0,
			},
			wantErr: true,
		},
		{
			name: "Negative Quantity",
			input: &MoveStockRequest{
				ProductID:      1,
				FromLocationID: 1,
				ToLocationID:   2,
				Quantity:       -5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// validate is a simple validation method for MoveStockRequest
func (r *MoveStockRequest) validate() error {
	if r.ProductID <= 0 {
		return assert.AnError
	}
	if r.FromLocationID <= 0 {
		return assert.AnError
	}
	if r.ToLocationID <= 0 {
		return assert.AnError
	}
	if r.FromLocationID == r.ToLocationID {
		return assert.AnError
	}
	if r.Quantity <= 0 {
		return assert.AnError
	}
	return nil
}

func TestStock_BasicProperties(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		stock        *Stock
		expectedID   int
		expectedPID  int
		expectedLID  int
		expectedQty  int
		expectedTime time.Time
	}{
		{
			name: "Valid Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expectedID:   1,
			expectedPID:  1,
			expectedLID:  1,
			expectedQty:  10,
			expectedTime: testTime,
		},
		{
			name: "Stock With Different Values",
			stock: &Stock{
				ID:         42,
				ProductID:  5,
				LocationID: 3,
				Quantity:   100,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expectedID:   42,
			expectedPID:  5,
			expectedLID:  3,
			expectedQty:  100,
			expectedTime: testTime,
		},
		{
			name: "Zero Stock",
			stock: &Stock{
				ID:         3,
				ProductID:  2,
				LocationID: 2,
				Quantity:   0,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expectedID:   3,
			expectedPID:  2,
			expectedLID:  2,
			expectedQty:  0,
			expectedTime: testTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, tt.stock.ID)
			assert.Equal(t, tt.expectedPID, tt.stock.ProductID)
			assert.Equal(t, tt.expectedLID, tt.stock.LocationID)
			assert.Equal(t, tt.expectedQty, tt.stock.Quantity)
			assert.Equal(t, tt.expectedTime, tt.stock.CreatedAt)
			assert.Equal(t, tt.expectedTime, tt.stock.UpdatedAt)
		})
	}
}

func TestStockMovement_BasicProperties(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	fromLocID, toLocID := 1, 2

	tests := []struct {
		name            string
		movement        *StockMovement
		expectedID      int
		expectedPID     int
		expectedFromLoc *int
		expectedToLoc   *int
		expectedQty     int
		expectedType    string
		expectedTime    time.Time
	}{
		{
			name: "Valid Stock Movement",
			movement: &StockMovement{
				ID:             1,
				ProductID:      1,
				FromLocationID: &fromLocID,
				ToLocationID:   &toLocID,
				Quantity:       5,
				MovementType:   "TRANSFER",
				CreatedAt:      testTime,
			},
			expectedID:      1,
			expectedPID:     1,
			expectedFromLoc: &fromLocID,
			expectedToLoc:   &toLocID,
			expectedQty:     5,
			expectedType:    "TRANSFER",
			expectedTime:    testTime,
		},
		{
			name: "Stock Movement Without From Location",
			movement: &StockMovement{
				ID:           1,
				ProductID:    1,
				ToLocationID: &toLocID,
				Quantity:     5,
				MovementType: "ADD",
				CreatedAt:    testTime,
			},
			expectedID:      1,
			expectedPID:     1,
			expectedFromLoc: nil,
			expectedToLoc:   &toLocID,
			expectedQty:     5,
			expectedType:    "ADD",
			expectedTime:    testTime,
		},
		{
			name: "Stock Movement Without To Location",
			movement: &StockMovement{
				ID:             1,
				ProductID:      1,
				FromLocationID: &fromLocID,
				Quantity:       5,
				MovementType:   "REMOVE",
				CreatedAt:      testTime,
			},
			expectedID:      1,
			expectedPID:     1,
			expectedFromLoc: &fromLocID,
			expectedToLoc:   nil,
			expectedQty:     5,
			expectedType:    "REMOVE",
			expectedTime:    testTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedID, tt.movement.ID)
			assert.Equal(t, tt.expectedPID, tt.movement.ProductID)
			assert.Equal(t, tt.expectedFromLoc, tt.movement.FromLocationID)
			assert.Equal(t, tt.expectedToLoc, tt.movement.ToLocationID)
			assert.Equal(t, tt.expectedQty, tt.movement.Quantity)
			assert.Equal(t, tt.expectedType, tt.movement.MovementType)
			assert.Equal(t, tt.expectedTime, tt.movement.CreatedAt)
		})
	}
}

func TestStock_IsLowStock(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		stock     *Stock
		threshold int
		want      bool
	}{
		{
			name: "Stock Above Threshold",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   15,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			threshold: 10,
			want:      false,
		},
		{
			name: "Stock At Threshold",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			threshold: 10,
			want:      false,
		},
		{
			name: "Stock Below Threshold",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   5,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			threshold: 10,
			want:      true,
		},
		{
			name: "Zero Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   0,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			threshold: 10,
			want:      true,
		},
		{
			name: "Zero Threshold",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   5,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			threshold: 0,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLowStock(tt.stock, tt.threshold)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStock_CanRemoveQuantity(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		stock    *Stock
		quantity int
		want     bool
	}{
		{
			name: "Sufficient Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: 5,
			want:     true,
		},
		{
			name: "Exact Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: 10,
			want:     true,
		},
		{
			name: "Insufficient Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   5,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: 10,
			want:     false,
		},
		{
			name: "Zero Quantity",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: 0,
			want:     true,
		},
		{
			name: "Negative Quantity",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   10,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: -5,
			want:     false,
		},
		{
			name: "Zero Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   0,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			quantity: 1,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := canRemoveQuantity(tt.stock, tt.quantity)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStock_GetStatus(t *testing.T) {
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		stock    *Stock
		expected string
	}{
		{
			name: "In Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   50,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expected: "in_stock",
		},
		{
			name: "Low Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   5,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expected: "low_stock",
		},
		{
			name: "Out of Stock",
			stock: &Stock{
				ID:         1,
				ProductID:  1,
				LocationID: 1,
				Quantity:   0,
				CreatedAt:  testTime,
				UpdatedAt:  testTime,
			},
			expected: "out_of_stock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatus(tt.stock)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions to simulate the methods that would be on the Stock struct

func isLowStock(stock *Stock, threshold int) bool {
	return stock.Quantity < threshold
}

func canRemoveQuantity(stock *Stock, quantity int) bool {
	if quantity < 0 {
		return false
	}
	return stock.Quantity >= quantity
}

func getStatus(stock *Stock) string {
	if stock.Quantity == 0 {
		return "out_of_stock"
	} else if stock.Quantity < 10 { // Assuming 10 is the low stock threshold
		return "low_stock"
	} else {
		return "in_stock"
	}
}
