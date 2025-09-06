// Package models provides data structures for the inventory management system.
// It defines the core entities such as products, locations, stock, and stock movements
// that are used throughout the application.
package models

import (
	"time"
)

// Stock represents the quantity of a specific product at a specific location.
// It tracks the current inventory levels and includes timestamps for creation and last update.
type Stock struct {
	ID         int       `json:"id" db:"id"`
	ProductID  int       `json:"product_id" db:"product_id"`
	LocationID int       `json:"location_id" db:"location_id"`
	Quantity   int       `json:"quantity" db:"quantity"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// StockMovement represents a movement of stock from one location to another.
// It tracks the product, source and destination locations, quantity moved, and movement type.
type StockMovement struct {
	ID             int       `json:"id" db:"id"`
	ProductID      int       `json:"product_id" db:"product_id"`
	FromLocationID *int      `json:"from_location_id" db:"from_location_id"`
	ToLocationID   *int      `json:"to_location_id" db:"to_location_id"`
	Quantity       int       `json:"quantity" db:"quantity"`
	MovementType   string    `json:"movement_type" db:"movement_type"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// AddStockRequest represents the data needed to add stock to a location.
// It contains the product ID, location ID, and quantity to add.
type AddStockRequest struct {
	ProductID  int `json:"product_id" validate:"required"`
	LocationID int `json:"location_id" validate:"required"`
	Quantity   int `json:"quantity" validate:"required,min=1"`
}

// MoveStockRequest represents the data needed to move stock between locations.
// It contains the product ID, source location ID, destination location ID, and quantity to move.
type MoveStockRequest struct {
	ProductID      int `json:"product_id" validate:"required"`
	FromLocationID int `json:"from_location_id" validate:"required"`
	ToLocationID   int `json:"to_location_id" validate:"required"`
	Quantity       int `json:"quantity" validate:"required,min=1"`
}
