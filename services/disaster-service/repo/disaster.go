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
	UpdateStatus(ctx context.Context, disasterID, status string) error
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

// UpdateStatus updates the status of a disaster by its ID.
func (r *InMemoryDisasterRepo) UpdateStatus(ctx context.Context, disasterID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if disaster, exists := r.data[disasterID]; exists {
		disaster.Status = status
		return nil
	}
	return nil
}
