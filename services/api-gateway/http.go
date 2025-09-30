package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cprakhar/relief-ops/services/api-gateway/handler"
)

type httpServer struct {
	addr string
}

// newHTTPServer creates and returns a new HTTP server.
func newHTTPServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

// run starts the HTTP server and listens for incoming requests.
func (s *httpServer) run(ctx context.Context) error {
	h := handler.NewHTTPUserHandler()
	srv := &http.Server{
		Addr:    s.addr,
		Handler: h,
	}

	// Listen for incoming requests in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Gracefully shutdown the server on context cancellation
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		log.Printf("Shutting down HTTP server on %s", s.addr)
		srv.Shutdown(shutdownCtx)
	}()

	// Wait for either an error or context cancellation
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}
}
