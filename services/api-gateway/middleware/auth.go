package middleware

import (
	"net/http"

	grpcclient "github.com/cprakhar/relief-ops/services/api-gateway/grpc_client"
	pb "github.com/cprakhar/relief-ops/shared/proto/user"
	"github.com/gin-gonic/gin"
)

const CookieName = "auth_token"

// JWTAuthMiddleware validates JWT tokens from cookies and sets user info in context.
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

	ctx.Next()
}

// AdminOnlyMiddleware ensures that the user has an admin role.
func AdminOnlyMiddleware(ctx *gin.Context) {
	role := ctx.GetString("role")
	if role != "admin" {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden: Admins only"})
		return
	}

	ctx.Next()
}
