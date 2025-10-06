package grpcclient

import (
	"fmt"

	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
	pb "github.com/cprakhar/relief-ops/shared/proto/disaster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type disasterServiceClient struct {
	Client pb.DisasterServiceClient
	conn   *grpc.ClientConn
}

// NewDisasterServiceClient creates and returns a new gRPC client for the Disaster Service.
func NewDisasterServiceClient() (*disasterServiceClient, error) {
	disasterServiceURL := fmt.Sprintf("disaster-service%s", env.GetString("DISASTER_GRPC_ADDR", ":9002"))

	dialOpts := traces.DialOptionsWithTracing()
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	conn, err := grpc.NewClient(disasterServiceURL, dialOpts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewDisasterServiceClient(conn)
	return &disasterServiceClient{Client: client, conn: conn}, nil
}

// Close closes the gRPC connection.
func (dsc *disasterServiceClient) Close() error {
	return dsc.conn.Close()
}
