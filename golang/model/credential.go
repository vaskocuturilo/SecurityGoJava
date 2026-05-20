package model

import (
	"errors"
	"strings"
)

type Credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	ErrEmailRequired    = errors.New("email required")
	ErrPasswordRequired = errors.New("password required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

func (c *Credential) Validate() error {
	if strings.TrimSpace(c.Email) == "" {
		return ErrEmailRequired
	}

	if strings.TrimSpace(c.Password) == "" {
		return ErrPasswordRequired
	}

	if len(c.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}
