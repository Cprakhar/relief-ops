package main

import (
	"context"
	"net"

	"github.com/cprakhar/relief-ops/services/disaster-service/handler"
	"github.com/cprakhar/relief-ops/services/disaster-service/service"
	"github.com/cprakhar/relief-ops/shared/messaging"
	"google.golang.org/grpc"
)

type gRPCServer struct {
	addr string
	svc  service.DisasterService
	kc   *messaging.KafkaClient
}

func newgRPCServer(addr string, svc service.DisasterService, kc *messaging.KafkaClient) *gRPCServer {
	return &gRPCServer{addr: addr, svc: svc, kc: kc}
}

func (s *gRPCServer) run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()
	handler.NewDisastergRPCHandler(srv, s.svc, s.kc)

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
