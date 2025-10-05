package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cprakhar/relief-ops/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var QueryTimeout = 5 * time.Second

type mongodbResourceRepo struct {
	db *mongo.Collection
}

type ResourceRepo interface {
	AddResources(ctx context.Context, resources []*types.Resource) error
	GetNearbyResources(ctx context.Context, lat, lon float64, radiusMeters int) ([]*types.Resource, error)
}

// NewResourceRepo creates a new instance of mongodbResourceRepo.
func NewResourceRepo(ctx context.Context, db *mongo.Collection) (ResourceRepo, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	// Create geospatial index on location field
	geoIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "location", Value: "2dsphere"}},
		Options: options.Index().SetName("location_2dsphere"),
	}

	// Create TTL index on created_at field to auto-expire documents after 30 days
	ttlIndexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "created_at", Value: 1}},
		Options: options.Index().
			SetExpireAfterSeconds(3600 * 24 * 30). // 30 days
			SetName("created_at_ttl"),
	}

	indexModel := []mongo.IndexModel{geoIndexModel, ttlIndexModel}
	_, err := db.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create indexes: %v", err)
	}

	return &mongodbResourceRepo{db: db}, nil
}

// AddResources adds multiple resources to the repository.
func (r *mongodbResourceRepo) AddResources(ctx context.Context, resources []*types.Resource) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	if len(resources) == 0 {
		return nil
	}

	now := time.Now()

	var operations []mongo.WriteModel
	for _, resource := range resources {
		filter := bson.M{
			"name":         resource.Name,
			"amenity_type": resource.AmenityType,
		}

		update := bson.M{
			"$set": bson.M{
				"name":         resource.Name,
				"amenity_type": resource.AmenityType,
				"location":     resource.Location,
				"updated_at":   now,
			},
			"$setOnInsert": bson.M{
				"created_at": now,
			},
		}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	bulkOpts := options.BulkWrite().SetOrdered(false)
	_, err := r.db.BulkWrite(ctx, operations, bulkOpts)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("Duplicate key error while adding resources: %v", err)
			return nil
		}
		return fmt.Errorf("bulk write failed: %w", err)
	}

	return nil
}

// GetNearbyResources retrieves resources within a certain radius (in meters) of given coordinates.
func (r *mongodbResourceRepo) GetNearbyResources(ctx context.Context, lat, lon float64, radiusMeters int) ([]*types.Resource, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat}, // GeoJSON format is [longitude, latitude]
				},
				"$maxDistance": radiusMeters,
			},
		},
	}

	findOpts := options.Find().
		SetLimit(100). // Limit to 100 results
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.db.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var resources []*types.Resource
	for cursor.Next(ctx) {
		var resource types.Resource
		if err := cursor.Decode(&resource); err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}
