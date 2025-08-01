package models

import (
	"time"
)

type Stock struct {
	ID         int       `json:"id" db:"id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	LocationID int       `json:"location_id" db:"location_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type StockMovement struct {
	ID             int       `json:"id" db:"id"`
	ProductID      int       `json:"product_id" db:"product_id"`
	FromLocationID *int      `json:"from_location_id" db:"from_location_id"`
	ToLocationID   *int      `json:"to_location_id" db:"to_location_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	MovementType   string    `json:"movement_type" db:"movement_type"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type AddStockRequest struct {
	ProductID  int `json:"product_id" validate:"required"`
	LocationID int `json:"location_id" validate:"required"`
	Quantity   int `json:"quantity" validate:"required,min=1"`
}

type MoveStockRequest struct {
	ProductID      int `json:"product_id" validate:"required"`
	FromLocationID int `json:"from_location_id" validate:"required"`
	ToLocationID   int `json:"to_location_id" validate:"required"`
	Quantity       int `json:"quantity" validate:"required,min=1"`
}
