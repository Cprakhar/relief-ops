package main

import (
	"context"
	"net/http"
	"time"

	"github.com/cprakhar/relief-ops/services/api-gateway/handler"
)

type httpServer struct {
	addr string
}

func newHTTPServer(addr string) *httpServer {
	return &httpServer{addr: addr}
}

func (s *httpServer) run(ctx context.Context) error {

	h := handler.NewHTTPUserHandler()

	srv := &http.Server{
		Addr:    s.addr,
		Handler: h,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(shutdownCtx)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return nil
	}

}
