package service

import (
	"context"
	"time"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/cprakhar/relief-ops/shared/util"
)

const (
	AdminRole       = "admin"
	ContributorRole = "contributor"
)

type JwtConfig struct {
	Secret string
	Expiry time.Duration
}

type userService struct {
	repo   repo.UserRepo
	jwtCfg *JwtConfig
}

// UserService defines the interface for user service operations.
type UserService interface {
	CreateUser(ctx context.Context, user *types.User) (string, error)
	Login(ctx context.Context, email, password string) (*types.User, string, error)
	GetUserByID(ctx context.Context, id string) (*types.User, error)
	GetAdmins(ctx context.Context) ([]*types.User, error)
}

// NewUserService creates a new instance of userService.
func NewUserService(r repo.UserRepo, secret string, expiry time.Duration) UserService {
	return &userService{repo: r, jwtCfg: &JwtConfig{Secret: secret, Expiry: expiry}}
}

// CreateUser creates a new user entry.
func (s *userService) CreateUser(ctx context.Context, user *types.User) (string, error) {
	// Hash the password before storing
	hashedPassword, err := util.EncryptPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hashedPassword

	return s.repo.Create(ctx, user)
}

// GetAdmins retrieves all users with the admin role.
func (s *userService) GetAdmins(ctx context.Context) ([]*types.User, error) {
	return s.repo.GetAllByRole(ctx, AdminRole)
}

// Login authenticates a user and returns a JWT token upon successful authentication.
func (s *userService) Login(ctx context.Context, email, password string) (*types.User, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if ok := util.ValidatePassword(user.Password, password); !ok {
		return nil, "", err
	}

	userDetails := &util.UserDetails{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
	}

	token, err := util.GenerateToken(userDetails, s.jwtCfg.Secret, s.jwtCfg.Expiry)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	return s.repo.GetByID(ctx, id)
}
