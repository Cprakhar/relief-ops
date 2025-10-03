package http

import (
	"log"
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pbu "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/cprakhar/relief-ops/shared/response"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const CookieName = "auth_token"

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// RegisterUserHandler handles user registration requests.
func RegisterUserHandler(ctx *gin.Context) {
	var req registerUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.JSONResponse{Error: err.Error()})
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
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response.JSONResponse{Data: gin.H{
		"user_id": pbRes.GetId(),
	}})
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginUserHandler handles user login requests.
func LoginUserHandler(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.JSONResponse{Error: err.Error()})
		return
	}

	userClient, err := grpcclient.NewUserServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Close()

	pbReq := &pbu.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	pbRes, err := userClient.Client.LoginUser(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	token := pbRes.GetToken()
	user := pbRes.GetUser()

	ctx.SetCookie(CookieName,
		token,
		3600*24*7, // 7 days
		"/",
		"",
		false,
		true,
	)
	ctx.SetSameSite(http.SameSiteStrictMode)

	responseData := struct {
		Token string    `json:"token"`
		User  *pbu.User `json:"user"`
	}{
		Token: token,
		User:  user,
	}

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: responseData})
}

// GetCurrentUserHandler retrieves the currently authenticated user's details.
func GetCurrentUserHandler(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, response.JSONResponse{Error: "Unauthorized"})
		return
	}

	userClient, err := grpcclient.NewUserServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Close()

	pbReq := &pbu.GetUserRequest{Id: userID}
	pbRes, err := userClient.Client.GetUser(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	oid, err := bson.ObjectIDFromHex(pbRes.GetId())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: "Invalid user ID"})
		return
	}

	user := &types.User{
		ID:        oid,
		Name:      pbRes.GetName(),
		Email:     pbRes.GetEmail(),
		Role:      pbRes.GetRole(),
		AvatarURL: pbRes.GetAvatarUrl(),
	}

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: user})
}

// LogoutUserHandler handles user logout by clearing the auth cookie.
func LogoutUserHandler(ctx *gin.Context) {
	ctx.SetCookie(CookieName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, response.JSONResponse{Data: "Logged out successfully"})
}
