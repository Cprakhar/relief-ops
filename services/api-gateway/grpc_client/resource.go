package grpcclient

import (
	"fmt"

	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
	pb "github.com/cprakhar/relief-ops/shared/proto/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type resourceServiceClient struct {
	Client pb.ResourceServiceClient
	conn   *grpc.ClientConn
}

// NewresourceServiceClient creates and returns a new gRPC client for the User Service.
func NewResourceServiceClient() (*resourceServiceClient, error) {
	resourceServiceURL := fmt.Sprintf("resource-service%s", env.GetString("RESOURCE_GRPC_ADDR", ":9001"))

	dialOpts := traces.DialOptionsWithTracing()
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(resourceServiceURL, dialOpts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewResourceServiceClient(conn)
	return &resourceServiceClient{Client: client, conn: conn}, nil
}

// Close closes the gRPC connection.
func (rsc *resourceServiceClient) Close() error {
	return rsc.conn.Close()
}
