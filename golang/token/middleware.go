package token

import (
	"context"
	"fmt"
	"golang/claims"
	"golang/model"
	"golang/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey = contextKey("user")

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		parts := strings.Fields(authHeader)

		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		u, err := VerifyAccessToken(parts[1])

		if err != nil {
			log.Printf("Token error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := putUserToContext(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

func putUserToContext(ctx context.Context, u model.User) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}
