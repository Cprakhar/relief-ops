package main

import (
	"context"
	"log"
	"net"

	"github.com/cprakhar/relief-ops/services/user-service/handler"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	svc  service.UserService
}

// newgRPCServer creates a new gRPC server instance.
func newgRPCServer(addr string, svc service.UserService) *gRPCServer {
	return &gRPCServer{addr: addr, svc: svc}
}

// run starts the gRPC server and listens for incoming requests.
func (s *gRPCServer) run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	// Create a new gRPC server
	srv := grpc.NewServer()
	handler.NewUsergRPCHandler(srv, s.svc)

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
		log.Println("Gracefully stopping gRPC server...")
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
