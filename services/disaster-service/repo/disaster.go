package repo

import (
	"context"
	"sync"

	"github.com/cprakhar/relief-ops/services/disaster-service/types"
)

type InMemoryDisasterRepo struct {
	mu   sync.RWMutex
	data map[string]*types.Disaster
}

type DisasterRepo interface {
	Create(ctx context.Context, disaster *types.Disaster) (string, error)
	Delete(ctx context.Context, disasterID string) error
}

func NewDisasterRepo() *InMemoryDisasterRepo {
	return &InMemoryDisasterRepo{
		data: make(map[string]*types.Disaster),
	}
}

func (r *InMemoryDisasterRepo) Create(ctx context.Context, disaster *types.Disaster) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	disaster.Status = "pending"
	r.data[disaster.ID] = disaster
	return disaster.ID, nil
}

func (r *InMemoryDisasterRepo) Delete(ctx context.Context, disasterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, disasterID)
	return nil
}
