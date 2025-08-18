// Package repository provides data access implementations for the inventory management system.
// It contains repositories that handle interactions with the database for products, locations,
// stock, and stock movements.
package repository

import (
	"context"
	"fmt"

	"cli-inventory/internal/db"
	"cli-inventory/internal/models"
)

// LocationRepository provides methods for interacting with location data in the database.
// It implements the LocationRepositoryInterface defined in the service package.
type LocationRepository struct {
	queries *db.Queries
}

// NewLocationRepository creates a new instance of LocationRepository with the provided database queries.
func NewLocationRepository(queries *db.Queries) *LocationRepository {
	return &LocationRepository{
		queries: queries,
	}
}

func (r *LocationRepository) Create(ctx context.Context, location *models.CreateLocationRequest) (*models.Location, error) {
	dbLocation, err := r.queries.CreateLocation(ctx, location.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt.Time,
	}, nil
}

func (r *LocationRepository) GetByName(ctx context.Context, name string) (*models.Location, error) {
	dbLocation, err := r.queries.GetLocationByName(ctx, name)
	if err != nil {
		// If no location is found, return nil instead of an error
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get location by name: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt.Time,
	}, nil
}

func (r *LocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	dbLocation, err := r.queries.GetLocationByID(ctx, int32(id))
	if err != nil {
		// If no location is found, return nil instead of an error
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get location by ID: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt.Time,
	}, nil
}

func (r *LocationRepository) List(ctx context.Context) ([]models.Location, error) {
	dbLocations, err := r.queries.ListLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}

	locations := make([]models.Location, len(dbLocations))
	for i, dbLocation := range dbLocations {
		locations[i] = models.Location{
			ID:        int(dbLocation.ID),
			Name:      dbLocation.Name,
			CreatedAt: dbLocation.CreatedAt.Time,
		}
	}

	return locations, nil
}
