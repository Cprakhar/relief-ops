package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	QueryTimeout = 5 * time.Second
)

type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"`
	AvatarURL string    `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	Role      string    `json:"role" bson:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// mongodbUserRepo is a simple in-memory implementation of UserRepo. (later switch with postgres)
type mongodbUserRepo struct {
	db *mongo.Collection
}

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Create(ctx context.Context, user *User) (string, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAllByRole(ctx context.Context, role string) ([]*User, error)
}

// NewUserRepo creates a new instance of inMemoryUserRepo.
func NewUserRepo(db *mongo.Collection) UserRepo {
	return &mongodbUserRepo{db: db}
}

// Create adds a new user to the repository.
func (r *mongodbUserRepo) Create(ctx context.Context, user *User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Role = "contributor" // default role

	res, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(string), nil
}

// GetByID retrieves a user by their ID.
func (r *mongodbUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user User
	err := r.db.FindOne(ctx, map[string]string{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email.
func (r *mongodbUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	var user User
	err := r.db.FindOne(ctx, map[string]string{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllByRole retrieves all users with the specified role.
func (r *mongodbUserRepo) GetAllByRole(ctx context.Context, role string) ([]*User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	cursor, err := r.db.Find(ctx, map[string]string{"role": role})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*User
	for cursor.Next(ctx) {
		var user User
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
