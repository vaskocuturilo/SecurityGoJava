package model

import "errors"

type Credential struct {
	UserName string `json:"userName"`
	Password []byte `json:"password"`
}

var (
	ErrUsernameRequired = errors.New("username required")
	ErrPasswordRequired = errors.New("password required")
	ErrInvalidInput     = errors.New("invalid input data")
	ErrAlreadyExists    = errors.New("event already exists")
)

func (c *Credential) Validate() error {
	if c.UserName == "" {
		return ErrUsernameRequired
	}

	if len(c.Password) == 0 {
		return ErrPasswordRequired
	}

	return nil
}
