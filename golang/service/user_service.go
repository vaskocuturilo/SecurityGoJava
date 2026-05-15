package service

import (
	"context"
	"golang/model"
	"golang/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SignUp(ctx context.Context, credential *model.Credential) error {
	if err := credential.Validate(); err != nil {
		return err
	}
	return s.repo.SignUp(ctx, credential)
}

func (s *UserService) Login(ctx context.Context, user *model.User) error {
	return s.repo.Login(ctx, user)
}

func (s *UserService) Refresh(ctx context.Context, request *model.RefreshRequest) error {
	return nil
}
