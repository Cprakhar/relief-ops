package repo

import (
	"context"
	"sync"

	"github.com/cprakhar/relief-ops/shared/types"
)

type InMemoryDisasterRepo struct {
	mu   sync.RWMutex
	data map[string]*types.Disaster
}

// DisasterRepo defines the interface for disaster repository operations.
type DisasterRepo interface {
	Create(ctx context.Context, disaster *types.Disaster) (string, error)
	Delete(ctx context.Context, disasterID string) error
}

// NewDisasterRepo creates a new instance of InMemoryDisasterRepo.
func NewDisasterRepo() *InMemoryDisasterRepo {
	return &InMemoryDisasterRepo{
		data: make(map[string]*types.Disaster),
	}
}

// Create adds a new disaster to the repository.
func (r *InMemoryDisasterRepo) Create(ctx context.Context, disaster *types.Disaster) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	disaster.Status = "pending"
	r.data[disaster.ID] = disaster
	return disaster.ID, nil
}

// Delete removes a disaster from the repository by its ID.
func (r *InMemoryDisasterRepo) Delete(ctx context.Context, disasterID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, disasterID)
	return nil
}
