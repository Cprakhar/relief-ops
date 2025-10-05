package http

import (
	"log"
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	"github.com/cprakhar/relief-ops/shared/env"
	pbu "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/cprakhar/relief-ops/shared/response"
	"github.com/cprakhar/relief-ops/shared/types"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const CookieName = "auth_token"

type registerUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role" binding:"required"`
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
		Role:     req.Role,
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

type OAuthProvidersConfig struct {
	Google OAuthProvider
	GitHub OAuthProvider
}

type OAuthProvider struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// InitOAuthProviders initializes OAuth providers using the provided configuration.
func InitOAuthProviders(oauthCfg *OAuthProvidersConfig) {
	google.New(
		oauthCfg.Google.ClientID,
		oauthCfg.Google.ClientSecret,
		oauthCfg.Google.RedirectURL,
		oauthCfg.Google.Scopes...,
	)
	github.New(
		oauthCfg.GitHub.ClientID,
		oauthCfg.GitHub.ClientSecret,
		oauthCfg.GitHub.RedirectURL,
		oauthCfg.GitHub.Scopes...,
	)
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

func OAuthSignInHandler(ctx *gin.Context) {
	provider := ctx.Query("provider")

	if provider == "" {
		ctx.JSON(http.StatusBadRequest, response.JSONResponse{Error: "Provider is required"})
		return
	}

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func OAuthCallbackHandler(ctx *gin.Context) {
	provider := ctx.Query("provider")
	role := ctx.Query("role")

	if role == "" {
		role = "user"
	}

	if provider == "" {
		ctx.JSON(http.StatusBadRequest, response.JSONResponse{Error: "Provider is required"})
		return
	}

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	userClient, err := grpcclient.NewUserServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer userClient.Close()

	pbReq := &pbu.OAuthSignInRequest{
		Email:     user.Email,
		Name:      user.Name,
		AvatarUrl: user.AvatarURL,
		Role:      role,
	}

	pbRes, err := userClient.Client.OAuthSignIn(ctx, pbReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.JSONResponse{Error: err.Error()})
		return
	}

	token := pbRes.GetToken()

	ctx.SetCookie(CookieName,
		token,
		3600*24*7, // 7 days
		"/",
		"",
		false,
		true,
	)
	ctx.SetSameSite(http.SameSiteStrictMode)

	webURL := env.GetString("WEB_URL", "http://localhost:3000")

	ctx.Redirect(http.StatusTemporaryRedirect, webURL+"/oauth-success?token="+token)
}
