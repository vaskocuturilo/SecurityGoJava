package repository

import (
	"context"
	"golang/model"
)

type UserRepository interface {
	SignUp(ctx context.Context, credential *model.Credential) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}
