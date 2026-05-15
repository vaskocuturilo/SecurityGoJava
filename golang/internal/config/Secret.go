package config

import (
	"golang/model"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

const issuer = "example.com"

func JWTSecret() []byte {
	return jwtSecretKey
}

func UserDB() []model.User {
	return usersDB
}

func GetIssuer() string {
	return issuer
}

func hashPassword(password string) []byte {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hash
}

func AccessTokenDuration() time.Duration {
	return time.Duration(time.Now().Add(15 * time.Minute).Unix())
}

var usersDB = []model.User{
	{
		Name:           "John",
		Email:          "john.doe@test.com",
		HashedPassword: hashPassword("john.doe.password"),
	},
	{
		Name:           "Jane",
		Email:          "jane.doe@test.com",
		HashedPassword: hashPassword("jane.doe.password"),
	},
}
