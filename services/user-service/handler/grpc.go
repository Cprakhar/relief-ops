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

func NewUsergRPCHandler(srv *grpc.Server, svc service.UserService) {
	handler := &gRPCHandler{svc: svc}
	pb.RegisterUserServiceServer(srv, handler)
}

func (h *gRPCHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	name := req.GetName()

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

	userID, err := h.svc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserID: userID,
	}, nil
}

func (h *gRPCHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	user, err := h.svc.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

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
