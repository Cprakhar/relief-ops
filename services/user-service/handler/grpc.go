package handler

import (
	"context"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/cprakhar/relief-ops/shared/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedUserServiceServer
	svc service.UserService
}

// NewUsergRPCHandler registers the gRPC handler for user service.
func NewUsergRPCHandler(srv *grpc.Server, svc service.UserService) {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterUserServiceServer(srv, handler)
}

// RegisterUser handles user registration.
func (h *gRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	name := req.GetName()
	
	// Check if user already exists
	if exists := h.svc.UserExists(ctx, email); exists {
		return nil, status.Errorf(409, "user with email %s already exists", email)
	}

	// Hash the password before storing
	hashedPassword, err := util.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	user := &repo.User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: hashedPassword,
	}

	// Create the user
	userID, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		Id: userID,
	}, nil
}

// LoginUser handles user login.
func (h *gRPCHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	// Fetch user by email
	user, err := h.svc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Validate password
	valid := util.ValidatePassword(user.Password, password)
	if !valid {
		return nil, status.Errorf(401, "invalid credentials")
	}

	// Generate a token (this is a placeholder, implement your own token generation logic)
	token, err := util.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginUserResponse{
		Id:    user.ID,
		Token: token,
	}, nil

}
