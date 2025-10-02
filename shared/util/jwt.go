package util

import (
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims with custom user claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// UserDetails holds the information about the user for whom the token is generated.
type UserDetails struct {
	UserID string
	Email  string
	Role   string
}

var (
	Iss = "relief-ops"
	Aud = "relief-ops-users"
	Sub = "user-authentication"
)

// GenerateToken creates a JWT token for the given user details.
func GenerateToken(user *UserDetails, secretKey string, expiry time.Duration) (string, error) {
	claims := &Claims{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Iss,
			Audience:  []string{"relief-ops-users"},
			Subject:   Sub,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// ParseToken validates the JWT token and returns the claims if valid.
func ParseToken(tokenStr, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !(ok && token.Valid) {
		return nil, fmt.Errorf("invalid token claims")
	}

	if err := validateClaims(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// validateClaims checks the standard claims of the token.
func validateClaims(claims *Claims) error {
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return fmt.Errorf("token has expired")
	}

	if claims.Issuer != Iss {
		return fmt.Errorf("invalid token issuer")
	}

	foundAudience := slices.Contains(claims.Audience, Aud)
	if !foundAudience {
		return fmt.Errorf("invalid token audience")
	}

	if claims.Subject != Sub {
		return fmt.Errorf("invalid token subject")
	}
	return nil
}
