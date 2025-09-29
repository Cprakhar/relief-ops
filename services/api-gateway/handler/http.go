package handler

import (
	"log"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/gin-gonic/gin"
)

func NewHTTPUserHandler() *gin.Engine {
	r := gin.Default()

	r.GET("/health", HealthCheckHandler)

	r.POST("/register", RegisterUserHandler)
	return r
}

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}

func RegisterUserHandler(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	client, err := grpcclient.NewUserServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	pbReq := &pb.RegisterUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	pbRes, err := client.Client.RegisterUser(ctx, pbReq)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"user_id": pbRes.GetUserID()})
}
