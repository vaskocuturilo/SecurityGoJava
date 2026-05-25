package service

import (
	"context"
	"golang/model"
)

type IUserService interface {
	SignUp(ctx context.Context, credential *model.Credential) error
	Login(ctx context.Context, email, password string) (*model.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}
