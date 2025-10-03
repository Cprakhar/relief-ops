package db

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDBConfig struct {
	URI        string
	Database   string
	Collection string
	Timeout    *time.Duration
	MaxIdle    *time.Duration
	MaxPool    *uint64
	MinPool    *uint64
}

func NewMongoDBClient(cfg *MongoDBConfig) (*mongo.Collection, error) {
	serverAPIOpts := options.ServerAPI(options.ServerAPIVersion1)

	clientOpts := &options.ClientOptions{
		Timeout:         cfg.Timeout,
		MaxConnIdleTime: cfg.MaxIdle,
		MaxPoolSize:     cfg.MaxPool,
		MinPoolSize:     cfg.MinPool,
		BSONOptions: &options.BSONOptions{
			UseJSONStructTags: true,
		},
		ServerAPIOptions: serverAPIOpts,
	}

	clientOpts = clientOpts.ApplyURI(cfg.URI)
	mongodb, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	collection := mongodb.Database(cfg.Database).Collection(cfg.Collection)
	return collection, nil
}

func PrimitiveToHex(id interface{}) (string, error) {
	oid, ok := id.(bson.ObjectID)
	if !ok {
		return "", fmt.Errorf("invalid ObjectID type")
	}
	return oid.Hex(), nil
}
