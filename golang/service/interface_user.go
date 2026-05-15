package service

import (
	"context"
	"golang/model"
)

type IUserService interface {
	SignUp(ctx context.Context, credential *model.Credential) error
	Login(ctx context.Context, credential *model.User) error
	Refresh(ctx context.Context, request *model.RefreshRequest) error
}
