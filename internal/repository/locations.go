package repository

import (
	"context"
	"fmt"

	"github.com/rodrigotoledo/cli-inventory/internal/db"
	"github.com/rodrigotoledo/cli-inventory/internal/models"
)

type LocationRepository struct {
	queries *db.Queries
}

func NewLocationRepository(queries *db.Queries) *LocationRepository {
	return &LocationRepository{
		queries: queries,
	}
}

func (r *LocationRepository) Create(ctx context.Context, location *models.CreateLocationRequest) (*models.Location, error) {
	params := db.CreateLocationParams{
		Name: location.Name,
	}

	dbLocation, err := r.queries.CreateLocation(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt,
	}, nil
}

func (r *LocationRepository) GetByName(ctx context.Context, name string) (*models.Location, error) {
	dbLocation, err := r.queries.GetLocationByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get location by name: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt,
	}, nil
}

func (r *LocationRepository) GetByID(ctx context.Context, id int) (*models.Location, error) {
	dbLocation, err := r.queries.GetLocationByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get location by ID: %w", err)
	}

	return &models.Location{
		ID:        int(dbLocation.ID),
		Name:      dbLocation.Name,
		CreatedAt: dbLocation.CreatedAt,
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
			CreatedAt: dbLocation.CreatedAt,
		}
	}

	return locations, nil
}
