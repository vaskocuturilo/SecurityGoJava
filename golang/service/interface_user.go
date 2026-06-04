package service

import (
	"context"
	"golang/model"
)

type IUserService interface {
	SignUp(ctx context.Context, credential *model.Credential) error
	Login(ctx context.Context, email, password string) (*model.User, error)
	Refresh(refreshToken string) (string, string, error)
	GenerateTokens(user *model.User) (string, string, error)
}
