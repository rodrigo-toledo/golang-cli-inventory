// Code generated helpers to map sqlc/db types to internal models.
// Centralizing mapping logic avoids repetitive pgtype handling and keeps repository methods concise.
package repository

import (
	"cli-inventory/internal/db"
	"cli-inventory/internal/models"
)

// mapDBProductToModel converts a db.Product (sqlc generated) to *models.Product.
// It safely handles nullable pgtypes coming from the database.
func mapDBProductToModel(dbProduct db.Product) *models.Product {
	// Description (pgtype.Text)
	descriptionStr := ""
	if dbProduct.Description.Valid {
		descriptionStr = dbProduct.Description.String
	}

	// Price (pgtype.Numeric -> float64)
	var priceFloat float64
	if dbProduct.Price.Valid {
		if val, err := dbProduct.Price.Value(); err == nil {
			if floatVal, ok := val.(float64); ok {
				priceFloat = floatVal
			}
		}
	}

	return &models.Product{
		ID:          int(dbProduct.ID),
		SKU:         dbProduct.Sku,
		Name:        dbProduct.Name,
		Description: descriptionStr,
		Price:       priceFloat,
		CreatedAt:   dbProduct.CreatedAt.Time,
	}
}

// mapDBProductsToModels converts a slice of db.Product to a slice of models.Product.
func mapDBProductsToModels(dbProducts []db.Product) []models.Product {
	products := make([]models.Product, len(dbProducts))
	for i, p := range dbProducts {
		mp := mapDBProductToModel(p)
		products[i] = *mp
	}
	return products
}
