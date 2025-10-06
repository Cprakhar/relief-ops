package main

import (
	"context"
	"net"

	"github.com/cprakhar/relief-ops/services/resource-service/handler"
	"github.com/cprakhar/relief-ops/services/resource-service/service"
	"github.com/cprakhar/relief-ops/shared/messaging"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	svc  service.ResourceService
	kc   *messaging.KafkaClient
}

// newgRPCServer creates a new gRPC server instance.
func newgRPCServer(addr string, svc service.ResourceService, kc *messaging.KafkaClient) *gRPCServer {
	return &gRPCServer{addr: addr, svc: svc, kc: kc}
}

// run starts the gRPC server and listens for incoming requests.
func (s *gRPCServer) run(ctx context.Context) error {
	logger := logs.L()

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	// Create a new gRPC server
	srv := grpc.NewServer(traces.WithTracingInterceptors()...)
	handler.NewResourcegRPCHandler(srv, s.svc)

	// Listen for incoming requests in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}()

	// Gracefully shutdown the server on context cancellation
	go func() {
		<-ctx.Done()
		logger.Info("Gracefully stopping gRPC server...")
		srv.GracefulStop()
	}()

	// Wait for either an error or context cancellation
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
