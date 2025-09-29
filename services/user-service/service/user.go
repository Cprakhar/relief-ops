package service

import (
	"context"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
)

type userService struct {
	repo repo.UserRepo
}

type UserService interface {
	CreateUser(ctx context.Context, user *repo.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*repo.User, error)
}

func NewUserService(r repo.UserRepo) *userService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(ctx context.Context, user *repo.User) (string, error) {
	return s.repo.Create(ctx, user)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*repo.User, error) {
	return s.repo.GetByEmail(ctx, email)
}
