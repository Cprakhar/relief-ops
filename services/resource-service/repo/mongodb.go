package repo

import (
	"context"
	"fmt"

	"github.com/cprakhar/relief-ops/shared/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Resource struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Location types.Coordinates `json:"location"`
}

type mongodbResourceRepo struct {
	db *mongo.Collection
}

type ResourceRepo interface {
	AddResources(ctx context.Context, resources []*Resource) error
}

// NewResourceRepo creates a new instance of mongodbResourceRepo.
func NewResourceRepo(db *mongo.Collection) *mongodbResourceRepo {
	return &mongodbResourceRepo{db: db}
}

// AddResources adds multiple resources to the repository.
func (r *mongodbResourceRepo) AddResources(ctx context.Context, resources []*Resource) error {
	res, err := r.db.InsertMany(ctx, resources)
	if err != nil {
		return err
	}
	if len(res.InsertedIDs) != len(resources) {
		return fmt.Errorf("some resources were not inserted")
	}
	if !res.Acknowledged {
		return fmt.Errorf("insertion not acknowledged")
	}

	return nil
}
