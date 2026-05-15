package repository

import (
	"context"
	"golang/model"
)

type UserRepository interface {
	SignUp(ctx context.Context, credential *model.Credential) error
	Login(ctx context.Context, credential *model.User) error
	Refresh(ctx context.Context, request *model.RefreshRequest) error
}
