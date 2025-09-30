package repo

import (
	"context"
	"sync"
	"time"
)

type User struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// inMemoryUserRepo is a simple in-memory implementation of UserRepo. (later switch with postgres)
type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Create(ctx context.Context, user *User) (string, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByRole(ctx context.Context, role string) ([]*User, error)
}

// NewUserRepo creates a new instance of inMemoryUserRepo.
func NewUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users: make(map[string]*User),
	}
}

// Create adds a new user to the repository.
func (r *inMemoryUserRepo) Create(ctx context.Context, user *User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return user.ID, nil
}

// GetByID retrieves a user by their ID.
func (r *inMemoryUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

// GetByEmail retrieves a user by their email.
func (r *inMemoryUserRepo) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, nil
}

// GetByRole retrieves all users with the specified role.
func (r *inMemoryUserRepo) GetByRole(ctx context.Context, role string) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var admins []*User
	for _, user := range r.users {
		if user.Role == role {
			admins = append(admins, user)
		}
	}
	return admins, nil
}
