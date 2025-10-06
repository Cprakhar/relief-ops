package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cprakhar/relief-ops/shared/db"
	"github.com/cprakhar/relief-ops/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	QueryTimeout = 5 * time.Second
	ErrNotFound  = fmt.Errorf("record not found")
)

type mongodbDisasterRepo struct {
	db *mongo.Collection
}

// DisasterRepo defines the interface for disaster repository operations.
type DisasterRepo interface {
	Create(ctx context.Context, disaster *types.Disaster) (string, error)
	GetByID(ctx context.Context, id string) (*types.Disaster, error)
	GetAll(ctx context.Context, status string) ([]*types.Disaster, error)
	Delete(ctx context.Context, disasterID string) error
	UpdateStatus(ctx context.Context, disasterID, status string) error
}

// NewMongodbDisasterRepo creates a new instance of mongodbDisasterRepo.
func NewMongodbDisasterRepo(ctx context.Context, db *mongo.Collection) (DisasterRepo, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	ttlIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(3600 * 24 * 30). // 30 days
			SetName("created_at_ttl"),
	}

	_, err := db.Indexes().CreateOne(ctx, ttlIndexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	return &mongodbDisasterRepo{db: db}, nil
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

	return db.PrimitiveToHex(res.InsertedID)
}

// Delete deletes a disaster entry by its ID.
func (r *mongodbDisasterRepo) Delete(ctx context.Context, disasterID string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	filter := bson.M{
		"_id": disasterID,
	}

	res := r.db.FindOneAndDelete(ctx, filter)
	switch res.Err() {
	case mongo.ErrNoDocuments:
		return ErrNotFound
	default:
		return res.Err()
	}
}

func (r *mongodbDisasterRepo) GetByID(ctx context.Context, id string) (*types.Disaster, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var disaster types.Disaster
	filter := bson.M{
		"_id": id,
	}

	err := r.db.FindOne(ctx, filter).Decode(&disaster)
	if err != nil {
		return nil, err
	}
	return &disaster, nil
}

// UpdateStatus updates the status of a disaster entry.
func (r *mongodbDisasterRepo) UpdateStatus(ctx context.Context, disasterID, status string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	filter := bson.M{
		"_id": disasterID,
	}

	res := r.db.FindOneAndUpdate(ctx, filter, update)
	switch res.Err() {
	case mongo.ErrNoDocuments:
		return ErrNotFound
	default:
		return res.Err()
	}
}

// GetAll retrieves all disaster entries, optionally filtered by status.
func (r *mongodbDisasterRepo) GetAll(ctx context.Context, status string) ([]*types.Disaster, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var disasters []*types.Disaster
	for cursor.Next(ctx) {
		var disaster types.Disaster
		if err := cursor.Decode(&disaster); err != nil {
			return nil, err
		}
		disasters = append(disasters, &disaster)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return disasters, nil
}
