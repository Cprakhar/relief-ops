package main

import (
	"context"
	"net"

	"github.com/cprakhar/relief-ops/services/user-service/handler"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	svc  service.UserService
}

func newgRPCServer(addr string, svc service.UserService) *gRPCServer {
	return &gRPCServer{addr: addr, svc: svc}
}

func (s *gRPCServer) run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	handler.NewUsergRPCHandler(srv, s.svc)

	errChan := make(chan error, 1)
	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			errChan <- err
		}
	}()

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
