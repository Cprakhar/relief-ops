package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/cprakhar/relief-ops/shared/db"
	types "github.com/cprakhar/relief-ops/shared/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	QueryTimeout        = 5 * time.Second
	ErrResourceConflict = fmt.Errorf("resource already exists")
	ErrNoResourcesFound = fmt.Errorf("no resources found")
)

// mongodbUserRepo is a simple in-memory implementation of UserRepo. (later switch with postgres)
type mongodbUserRepo struct {
	db *mongo.Collection
}

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Create(ctx context.Context, user *types.User) (string, error)
	GetByID(ctx context.Context, id string) (*types.User, error)
	GetByEmail(ctx context.Context, email string) (*types.User, error)
	GetAllByRole(ctx context.Context, role string) ([]*types.User, error)
}

// NewUserRepo creates a new instance of inMemoryUserRepo.
func NewUserRepo(ctx context.Context, db *mongo.Collection) (UserRepo, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	if _, err := db.Indexes().CreateOne(ctx, indexModel); err != nil {
		return nil, err
	}

	return &mongodbUserRepo{db: db}, nil
}

// Create adds a new user to the repository.
func (r *mongodbUserRepo) Create(ctx context.Context, user *types.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	res, err := r.db.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrResourceConflict
		}
		return "", err
	}

	return db.PrimitiveToHex(res.InsertedID)
}

// GetByID retrieves a user by their ID.
func (r *mongodbUserRepo) GetByID(ctx context.Context, id string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user types.User
	filter := bson.M{
		"_id": id,
	}

	err := r.db.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email.
func (r *mongodbUserRepo) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user types.User
	err := r.db.FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, ErrNoResourcesFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// GetAllByRole retrieves all users with the specified role.
func (r *mongodbUserRepo) GetAllByRole(ctx context.Context, role string) ([]*types.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	cursor, err := r.db.Find(ctx, map[string]string{"role": role})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*types.User
	for cursor.Next(ctx) {
		var user types.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
