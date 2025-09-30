package service

import (
	"context"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
)

const (
	AdminRole       = "admin"
	ContributorRole = "contributor"
)

type userService struct {
	repo repo.UserRepo
}

// UserService defines the interface for user service operations.
type UserService interface {
	CreateUser(ctx context.Context, user *repo.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*repo.User, error)
	UserExists(ctx context.Context, email string) bool
	GetAdmins(ctx context.Context) ([]*repo.User, error)
}

// NewUserService creates a new instance of userService.
func NewUserService(r repo.UserRepo) *userService {
	return &userService{repo: r}
}

// CreateUser creates a new user entry.
func (s *userService) CreateUser(ctx context.Context, user *repo.User) (string, error) {
	return s.repo.Create(ctx, user)
}

// GetUserByEmail retrieves a user by their email.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*repo.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

// UserExists checks if a user with the given email exists.
func (s *userService) UserExists(ctx context.Context, email string) bool {
	_, err := s.GetUserByEmail(ctx, email)
	return err == nil
}

// GetAdmins retrieves all users with the admin role.
func (s *userService) GetAdmins(ctx context.Context) ([]*repo.User, error) {
	return s.repo.GetByRole(ctx, AdminRole)
}
