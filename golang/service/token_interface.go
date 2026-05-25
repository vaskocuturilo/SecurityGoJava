package service

import "golang/model"

type ITokenManager interface {
	CreateAccessToken(user *model.User) (string, error)
	CreateRefreshToken(u model.User) (string, error)
	VerifyRefreshToken(refreshToken string) (model.User, error)
}
