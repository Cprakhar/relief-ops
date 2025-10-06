package main

import (
	"context"
	"net/http"
	"time"

	handlerhttp "github.com/cprakhar/relief-ops/services/api-gateway/handler/http"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
)

type httpServer struct {
	addr   string
	webURL string
}

// newHTTPServer creates and returns a new HTTP server.
func newHTTPServer(addr string, web string) *httpServer {
	return &httpServer{addr: addr, webURL: web}
}

// run starts the HTTP server and listens for incoming requests.
func (s *httpServer) run(ctx context.Context) error {
	logger := logs.L()

	h := handlerhttp.NewHttpHandler(s.webURL)
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
		logger.Info("Shutting down HTTP server")
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
