package service

import (
	"context"

	"github.com/cprakhar/relief-ops/services/disaster-service/repo"
	"github.com/cprakhar/relief-ops/shared/types"
)

type disasterService struct {
	repo repo.DisasterRepo
}

// DisasterService defines the interface for disaster service operations.
type DisasterService interface {
	CreateDisaster(ctx context.Context, disaster *types.Disaster) (string, error)
	DeleteDisaster(ctx context.Context, disasterID string) error
	GetDisaster(ctx context.Context, disasterID string) (*types.Disaster, error)
	GetAllDisasters(ctx context.Context, status string) ([]*types.Disaster, error)
	UpdateStatus(ctx context.Context, disasterID, status string) error
}

// NewDisasterService creates a new instance of disasterService.
func NewDisasterService(r repo.DisasterRepo) *disasterService {
	return &disasterService{repo: r}
}

// CreateDisaster creates a new disaster entry.
func (s *disasterService) CreateDisaster(ctx context.Context, disaster *types.Disaster) (string, error) {
	return s.repo.Create(ctx, disaster)
}

// DeleteDisaster deletes a disaster entry by its ID.
func (s *disasterService) DeleteDisaster(ctx context.Context, disasterID string) error {
	return s.repo.Delete(ctx, disasterID)
}

// UpdateStatus updates the status of a disaster entry.
func (s *disasterService) UpdateStatus(ctx context.Context, disasterID, status string) error {
	return s.repo.UpdateStatus(ctx, disasterID, status)
}

// GetDisaster retrieves a disaster entry by its ID.
func (s *disasterService) GetDisaster(ctx context.Context, disasterID string) (*types.Disaster, error) {
	return s.repo.GetByID(ctx, disasterID)
}

// GetAllDisasters retrieves all disaster entries, optionally filtered by status.
func (s *disasterService) GetAllDisasters(ctx context.Context, status string) ([]*types.Disaster, error) {
	return s.repo.GetAll(ctx, status)
}
