package grpcclient

import (
	"fmt"

	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type userServiceClient struct {
	Client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserServiceClient creates and returns a new gRPC client for the User Service.
func NewUserServiceClient() (*userServiceClient, error) {
	userServiceURL := fmt.Sprintf("user-service%s", env.GetString("USER_GRPC_ADDR", ":9001"))

	dialOpts := traces.DialOptionsWithTracing()
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(userServiceURL, dialOpts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &userServiceClient{Client: client, conn: conn}, nil
}

// Close closes the gRPC connection.
func (usc *userServiceClient) Close() error {
	return usc.conn.Close()
}
