package handler

import (
	"context"

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

// NewUsergRPCHandler registers the gRPC handler for user service.
func NewUsergRPCHandler(srv *grpc.Server, svc service.UserService, s string) {
	handler := &gRPCHandler{
		svc:       svc,
		jwtSecret: s,
	}
	pb.RegisterUserServiceServer(srv, handler)
}

// RegisterUser handles user registration.
func (h *gRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	name := req.GetName()

	user := &types.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	// Create the user
	userID, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.RegisterUserResponse{
		Id: userID,
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
