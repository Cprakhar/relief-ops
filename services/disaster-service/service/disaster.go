package service

import (
	"context"

	"github.com/cprakhar/relief-ops/services/disaster-service/repo"
	"github.com/cprakhar/relief-ops/services/disaster-service/types"
)

type disasterService struct {
	repo repo.DisasterRepo
}

type DisasterService interface {
	CreateDisaster(ctx context.Context, disaster *types.Disaster) (string, error)
	DeleteDisaster(ctx context.Context, disasterID string) error
}

func NewDisasterService(r repo.DisasterRepo) *disasterService {
	return &disasterService{repo: r}
}

func (s *disasterService) CreateDisaster(ctx context.Context, disaster *types.Disaster) (string, error) {
	return s.repo.Create(ctx, disaster)
}

func (s *disasterService) DeleteDisaster(ctx context.Context, disasterID string) error {
	return s.repo.Delete(ctx, disasterID)
}
