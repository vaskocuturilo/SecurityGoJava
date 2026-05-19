package service

import (
	"context"
	"golang/model"
	"golang/repository"

	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return &model.User{}, model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if err != nil {
		return &model.User{}, model.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) Refresh(ctx context.Context, request *model.RefreshRequest) error {
	return nil
}
