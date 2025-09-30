package handler

import (
	"log"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pbd "github.com/cprakhar/relief-ops/shared/proto/disaster"
	pbu "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/gin-gonic/gin"
)

func NewHTTPUserHandler() *gin.Engine {
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", HealthCheckHandler)

	// User endpoints
	r.POST("/users/register", RegisterUserHandler)

	// Disaster endpoints
	r.POST("/disasters", ReportDisasterHandler)
	return r
}

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// HealthCheckHandler responds with a simple status message.
func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

// RegisterUserHandler handles user registration requests.
func RegisterUserHandler(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userClient, err := grpcclient.NewUserServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Close()

	pbReq := &pbu.RegisterUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	pbRes, err := userClient.Client.RegisterUser(ctx, pbReq)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"user_id": pbRes.GetId()})
}

type reportDisasterRequest struct {
	types.Disaster
}

// ReportDisasterHandler handles disaster reporting requests.
func ReportDisasterHandler(ctx *gin.Context) {

	// Get the user ID from the context (set by authentication middleware)
	userID := "some-user-id"

	var req reportDisasterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	disasterClient, err := grpcclient.NewDisasterServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer disasterClient.Close()

	pbReq := &pbd.ReportDisasterRequest{
		Title:         req.Title,
		Description:   req.Description,
		Tags:          req.Tags,
		Location:      &pbd.Coordinates{Latitude: req.Location.Latitude, Longitude: req.Location.Longitude},
		ContributorID: userID,
		ImageURLs:     req.ImageURLs,
	}

	pbRes, err := disasterClient.Client.ReportDisaster(ctx, pbReq)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"disaster_id": pbRes.GetId()})
}
