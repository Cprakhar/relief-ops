package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cprakhar/relief-ops/services/api-gateway/handler/http"
	"github.com/cprakhar/relief-ops/shared/env"
)

var (
	addr   = env.GetString("API_GATEWAY_ADDR", ":8080")
	webURL = env.GetString("WEB_URL", "http://localhost:3000")

	// Google OAuth configuration
	googleClientID     = env.GetString("GOOGLE_CLIENT_ID", "")
	googleClientSecret = env.GetString("GOOGLE_CLIENT_SECRET", "")
	googleRedirectURL  = env.GetString("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/auth/callback?provider=google")

	// GitHub OAuth configuration
	githubClientID     = env.GetString("GITHUB_CLIENT_ID", "")
	githubClientSecret = env.GetString("GITHUB_CLIENT_SECRET", "")
	githubRedirectURL  = env.GetString("GITHUB_REDIRECT_URL", "http://localhost:8080/api/auth/callback?provider=github")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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
		log.Printf("API Gateway running on %s", addr)
		if err := httpServer.run(ctx); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	<-done
	log.Printf("API Gateway stopped")
}
