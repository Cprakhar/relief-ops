package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/api-gateway/handler/http"
	"github.com/cprakhar/relief-ops/shared/env"
	"github.com/cprakhar/relief-ops/shared/observe/logs"
	"github.com/cprakhar/relief-ops/shared/observe/traces"
)

var (
	addr        = env.GetString("API_GATEWAY_ADDR", ":8080")
	webURL      = env.GetString("WEB_URL", "http://localhost:3000")
	environment = env.GetString("ENVIRONMENT", "development")
	// OAuth configuration
	// Google OAuth configuration
	googleClientID     = env.GetString("GOOGLE_CLIENT_ID", "")
	googleClientSecret = env.GetString("GOOGLE_CLIENT_SECRET", "")
	googleRedirectURL  = env.GetString("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/callback?provider=google")
	// GitHub OAuth configuration
	githubClientID     = env.GetString("GITHUB_CLIENT_ID", "")
	githubClientSecret = env.GetString("GITHUB_CLIENT_SECRET", "")
	githubRedirectURL  = env.GetString("GITHUB_REDIRECT_URL", "http://localhost:8080/api/auth/callback?provider=github")

	// OTLP configuration
	otlpEndpoint = env.GetString("OTLP_ENDPOINT", "otel-collector:4317")
	otlpInsecure = env.GetBool("OTLP_INSECURE", true)
)

func main() {
	// Initialize logger
	logger, err := logs.Init("api-gateway")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Errorw("Failed to sync logger", "error", err)
		}
	}()
	logger.Info("Logger initialized")

	// Set up context with signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Initialize tracing
	tracerCfg := &traces.TracerConfig{
		ServiceName:      "api-gateway",
		Environment:      environment,
		Secure:           !otlpInsecure,
		ExporterEndpoint: otlpEndpoint,
	}

	shutdown, err := traces.InitTrace(ctx, tracerCfg)
	if err != nil {
		logger.Fatalw("Failed to initialize tracing", "error", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			logger.Fatalw("Failed to shutdown tracer provider", "error", err)
		}
	}()
	logger.Info("Tracing initialized")

	// OAuth configuration
	oauthCfg := &http.OAuthProvidersConfig{
		Google: http.OAuthProvider{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  googleRedirectURL,
			Scopes:       []string{"email", "profile"},
		},
		GitHub: http.OAuthProvider{
			ClientID:     githubClientID,
			ClientSecret: githubClientSecret,
			RedirectURL:  githubRedirectURL,
			Scopes:       []string{"user:email"},
		},
	}

	// Initialize OAuth providers
	http.InitOAuthProviders(oauthCfg)

	// Start HTTP server
	httpServer := newHTTPServer(addr, webURL)

	done := make(chan struct{})
	go func() {
		defer close(done)
		logger.Infow("API Gateway running", "addr", addr)
		if err := httpServer.run(ctx); err != nil {
			logger.Errorw("HTTP server error", "error", err)
		}
	}()

	<-ctx.Done()
	<-done
	logger.Info("API Gateway stopped")
}
