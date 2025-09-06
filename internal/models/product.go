package models

import (
	"time"
)

type Product struct {
	ID          int       `json:"id" db:"id"`
	SKU         string    `json:"sku" db:"sku" validate:"required"`
	Name        string    `json:"name" db:"name" validate:"required"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type CreateProductRequest struct {
	SKU         string  `json:"sku" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
