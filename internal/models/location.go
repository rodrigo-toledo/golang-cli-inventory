package models

import (
	"time"
)

type Location struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" validate:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateLocationRequest struct {
	Name string `json:"name" validate:"required"`
}
