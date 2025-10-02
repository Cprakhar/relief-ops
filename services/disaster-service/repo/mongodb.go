package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cprakhar/relief-ops/shared/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	QueryTimeout = 5 * time.Second
)

type mongodbDisasterRepo struct {
	db *mongo.Collection
}

// DisasterRepo defines the interface for disaster repository operations.
type DisasterRepo interface {
	Create(ctx context.Context, disaster *types.Disaster) (string, error)
	Delete(ctx context.Context, disasterID string) error
	UpdateStatus(ctx context.Context, disasterID, status string) error
}

// NewMongodbDisasterRepo creates a new instance of mongodbDisasterRepo.
func NewMongodbDisasterRepo(db *mongo.Collection) *mongodbDisasterRepo {
	return &mongodbDisasterRepo{db: db}
}

// Create creates a new disaster entry.
func (r *mongodbDisasterRepo) Create(ctx context.Context, disaster *types.Disaster) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	disaster.Status = "pending"
	disaster.CreatedAt = time.Now()
	disaster.UpdatedAt = time.Now()

	res, err := r.db.InsertOne(ctx, disaster)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

// Delete deletes a disaster entry by its ID.
func (r *mongodbDisasterRepo) Delete(ctx context.Context, disasterID string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	res := r.db.FindOneAndDelete(ctx, map[string]string{"_id": disasterID})
	switch res.Err() {
	case mongo.ErrNoDocuments:
		return fmt.Errorf("record not found")
	default:
		return res.Err()
	}
}

// UpdateStatus updates the status of a disaster entry.
func (r *mongodbDisasterRepo) UpdateStatus(ctx context.Context, disasterID, status string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	update := map[string]any{
		"$set": map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	res := r.db.FindOneAndUpdate(ctx, map[string]string{"_id": disasterID}, update)
	switch res.Err() {
	case mongo.ErrNoDocuments:
		return fmt.Errorf("record not found")
	default:
		return res.Err()
	}
}
