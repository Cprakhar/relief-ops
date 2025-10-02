package http

import (
	"strings"
	"time"

	"github.com/cprakhar/relief-ops/services/api-gateway/middleware"
	"github.com/cprakhar/relief-ops/shared/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewHttpHandler sets up the HTTP routes and returns a Gin engine.
func NewHttpHandler(webURLs string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  false,
		AllowOrigins:     strings.Split(webURLs, ","),
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
	}))

	// Health check endpoint
	r.GET("/health", HealthCheckHandler)

	// Admin endpoints
	r.POST("/admin/review/:id", middleware.JWTAuthMiddleware, ReviewDisasterHandler)

	// User endpoints
	r.POST("/users/register", RegisterUserHandler)
	r.POST("/users/login", LoginUserHandler)

	// Disaster endpoints
	r.POST("/disasters", middleware.JWTAuthMiddleware, ReportDisasterHandler)
	r.GET("/disasters", GetAllDisastersHandler)
	r.GET("/disasters/:id", GetDisasterByIDHandler)
	return r
}

// HealthCheckHandler responds with a simple status message.
func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(200, response.JSONResponse{Data: gin.H{"status": "ok"}})
}
