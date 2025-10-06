package handler

import (
	"context"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/cprakhar/relief-ops/shared/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedUserServiceServer
	svc       service.UserService
	jwtSecret string
}

// GrpcHandler defines the gRPC handler interface for user service.
type GrpcHandler interface {
	OAuthSignIn(ctx context.Context, req *pb.OAuthSignInRequest) (*pb.LoginUserResponse, error)
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error)
	GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)
	ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error)
}

// NewUsergRPCHandler registers the gRPC handler for user service.
func NewUsergRPCHandler(srv *grpc.Server, svc service.UserService, s string) {
	handler := &gRPCHandler{
		svc:       svc,
		jwtSecret: s,
	}
	pb.RegisterUserServiceServer(srv, handler)
}

// OAuthSignIn handles user sign-in via OAuth providers.
func (h *gRPCHandler) OAuthSignIn(ctx context.Context, req *pb.OAuthSignInRequest) (*pb.LoginUserResponse, error) {
	email := req.GetEmail()
	name := req.GetName()
	avatarURL := req.GetAvatarUrl()
	role := req.GetRole()

	user := &types.User{
		Name:      name,
		Email:     email,
		Role:      role,
		AvatarURL: avatarURL,
	}

	userID, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		if err == repo.ErrResourceConflict {
			// User already exists, fetch the existing user
			existingUser, err := h.svc.GetUserByEmail(ctx, user.Email)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to get existing user: %v", err)
			}
			user = existingUser
		}
	}

	token, err := h.svc.OAuthSignIn(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to sign in user: %v", err)
	}

	return &pb.LoginUserResponse{
		Token: token,
		User: &pb.User{
			Id:        userID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			AvatarUrl: user.AvatarURL,
		},
	}, nil
}

// RegisterUser handles user registration.
func (h *gRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	name := req.GetName()
	role := req.GetRole()

	user := &types.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Create the user
	userID, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.RegisterUserResponse{
		Id:   userID,
		Role: role,
	}, nil
}

// LoginUser handles user login.
func (h *gRPCHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	user, token, err := h.svc.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginUserResponse{
		Token: token,
		User: &pb.User{
			Id:    user.ID.Hex(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

// GetUser retrieves a user by their ID.
func (h *gRPCHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	userID := req.GetId()

	user, err := h.svc.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.User{
		Id:        user.ID.Hex(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarUrl: user.AvatarURL,
	}, nil
}

// ValidateToken checks the validity of a JWT token and returns the associated user details.
func (h *gRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	tokenStr := req.GetToken()

	userDetails, err := util.ParseToken(tokenStr, h.jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	return &pb.ValidateTokenResponse{
		User: &pb.User{
			Id:    userDetails.UserID,
			Email: userDetails.Email,
			Role:  userDetails.Role,
		},
	}, nil
}
