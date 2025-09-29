package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/shared/env"
)

var (
	addr = env.GetString("API_GATEWAY_ADDR", ":8080")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	httpServer := newHTTPServer(addr)

	done := make(chan struct{})
	go func() {
		defer close(done)
		log.Printf("API Gateway running on %s", addr)
		if err := httpServer.run(ctx); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	<-done
	log.Printf("API Gateway stopped")
}
