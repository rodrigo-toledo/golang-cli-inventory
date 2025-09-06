// Package models provides data structures for the inventory management system.
// It defines the core entities such as products, locations, stock, and stock movements
// that are used throughout the application.
package models

import (
	"time"
)

// Location represents a physical location where inventory is stored.
// It contains information about the location including its name and creation timestamp.
type Location struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateLocationRequest represents the data needed to create a new location.
// It contains the name of the location to be created.
type CreateLocationRequest struct {
	Name string `json:"name" validate:"required"`
}
