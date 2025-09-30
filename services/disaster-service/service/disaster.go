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
