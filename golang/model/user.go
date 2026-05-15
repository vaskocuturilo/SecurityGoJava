package model

import "github.com/google/uuid"

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"hashed-password"`
}

func NewUser(username, email string, hashedPassword []byte) *User {
	return &User{
		ID:             uuid.New().String(),
		Name:           username,
		Email:          email,
		HashedPassword: hashedPassword,
	}
}
