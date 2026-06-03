package config

import (
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const issuer = "example.com"

var (
	jwtSecretKey []byte
	once         sync.Once
)

func JWTSecret() []byte {
	once.Do(func() {
		key := os.Getenv("JWT_SECRET_KEY")
		if key == "" {
			panic("JWT_SECRET_KEY is not set in environment variables")
		}
		jwtSecretKey = []byte(key)
	})

	return jwtSecretKey
}

func GetIssuer() string {
	return issuer
}

func HashPassword(password string) ([]byte, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash, nil
}

func AccessTokenDuration() time.Duration {
	return 15 * time.Minute
}
