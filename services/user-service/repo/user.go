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

type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*User
}

type UserRepo interface {
	Create(ctx context.Context, user *User) (string, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

func NewUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{
		users: make(map[string]*User),
	}
}

func (r *inMemoryUserRepo) Create(ctx context.Context, user *User) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
	return user.ID, nil
}

func (r *inMemoryUserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

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
