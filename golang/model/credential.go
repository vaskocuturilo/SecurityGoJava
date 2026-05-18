package model

import (
	"errors"
	"strings"
)

type Credential struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrUsernameRequired = errors.New("username required")
	ErrPasswordRequired = errors.New("password required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

func (c *Credential) Validate() error {
	if strings.TrimSpace(c.UserName) == "" {
		return ErrUsernameRequired
	}

	if strings.TrimSpace(c.Password) == "" {
		return ErrPasswordRequired
	}

	if len(c.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
