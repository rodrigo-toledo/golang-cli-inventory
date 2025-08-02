// Package models provides data structures for the inventory management system.
// It defines the core entities such as products, locations, stock, and stock movements
// that are used throughout the application.
package models

import (
	"time"
)

// Product represents a product in the inventory system.
// It contains all the information about a product including its SKU, name,
// description, price, and creation timestamp.
type Product struct {
	ID          int       `json:"id" db:"id"`
	SKU         string    `json:"sku" db:"sku" validate:"required"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// CreateProductRequest represents the data needed to create a new product.
// It contains the SKU, name, description, and price of the product to be created.
type CreateProductRequest struct {
	SKU         string  `json:"sku" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
