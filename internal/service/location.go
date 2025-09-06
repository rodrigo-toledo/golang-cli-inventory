package service

import (
	"context"
	"fmt"

	"github.com/rodrigotoledo/cli-inventory/internal/models"
	"github.com/rodrigotoledo/cli-inventory/internal/repository"
)

type LocationService struct {
	repo *repository.LocationRepository
}

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
