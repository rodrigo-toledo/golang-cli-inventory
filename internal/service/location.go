// Package service provides business logic implementations for the inventory management system.
// It contains services that handle the core functionality such as product management,
// stock management, and location management.
package service

import (
	"context"
	"fmt"

	"cli-inventory/internal/models"
	"cli-inventory/internal/repository"
)

// LocationService provides methods for managing locations in the inventory system.
// It handles operations such as creating locations, retrieving location information,
// and listing all locations.
type LocationService struct {
	repo *repository.LocationRepository
}

// NewLocationService creates a new instance of LocationService with the provided location repository.
func NewLocationService(repo *repository.LocationRepository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) CreateLocation(ctx context.Context, req *models.CreateLocationRequest) (*models.Location, error) {
	// Check if location with this name already exists
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("location with name %s already exists", req.Name)
	}

	// Create the location
	location, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

func (s *LocationService) GetLocationByName(ctx context.Context, name string) (*models.Location, error) {
	location, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	return location, nil
}

func (s *LocationService) ListLocations(ctx context.Context) ([]models.Location, error) {
	locations, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations: %w", err)
	}
	return locations, nil
}
