package token

import (
	"golang/internal/config"
	"golang/model"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateAccessToken(t *testing.T) {
	mockUser := model.User{
		ID:    "f47ac10b-58cc-4372-a567",
		Name:  "John",
		Email: "Doe@doe.com",
	}

	tokenString, err := CreateAccessToken(mockUser)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret(), nil
	})

	if err != nil || !token.Valid {
		t.Errorf("Token is not valid: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Could not parse claims")
	}

	if claims["sub"] != mockUser.ID {
		t.Errorf("Expected sub %s, got %s", mockUser.ID, claims["sub"])
	}
}
