package token

import (
	"fmt"
	"golang/claims"
	"golang/model"
	"golang/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		u, err := VerifyAccessToken(parts[1])

		if err != nil {
			log.Printf("Token error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user", u)
		c.Next()
	}
}

func VerifyAccessToken(accessToken string) (model.User, error) {
	c := &claims.UserClaims{}

	_, err := jwt.ParseWithClaims(
		accessToken,
		c,
		utils.KeyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		return model.User{}, fmt.Errorf("token validation failed: %w", err)
	}

	return model.User{
		Username: c.Name,
		Email:    c.Email,
	}, nil
}
