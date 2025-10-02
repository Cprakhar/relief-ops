package middleware

import (
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/gin-gonic/gin"
)

const (
	CookieName = "auth_token"
)

func JWTAuthMiddleware(ctx *gin.Context) {
	cookie, err := ctx.Cookie(CookieName)
	if err != nil || cookie == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userClient, err := grpcclient.NewUserServiceClient()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer userClient.Close()

	pbReq := &pb.ValidateTokenRequest{Token: cookie}
	pbRes, err := userClient.Client.ValidateToken(ctx, pbReq)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := pbRes.GetUser()
	if user == nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	ctx.Set("user_id", user.GetId())
	ctx.Set("role", user.GetRole())
	ctx.Set("email", user.GetEmail())

	ctx.Next()
}
