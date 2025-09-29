package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
	"github.com/cprakhar/relief-ops/services/user-service/service"
	"github.com/cprakhar/relief-ops/shared/env"
)

var (
	addr = env.GetString("USER_GRPC_ADDR", ":9001")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userRepo := repo.NewUserRepo()
	userService := service.NewUserService(userRepo)

	gRPCServer := newgRPCServer(addr, userService)

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Printf("Server running on %s", addr)
		if err := gRPCServer.run(ctx); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()
	<-ctx.Done()
	<-done
}
