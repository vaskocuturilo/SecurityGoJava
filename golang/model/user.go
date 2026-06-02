package model

import (
	"errors"
)

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	HashedPassword string `json:"-"`
}

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAlreadyExists      = errors.New("user already exists")
)
