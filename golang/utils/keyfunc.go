package utils

import (
	"fmt"
	"golang/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

func KeyFunc(t *jwt.Token) (interface{}, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}

	return config.JWTSecret(), nil
}
