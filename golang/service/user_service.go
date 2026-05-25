package service

import (
	"context"
	"golang/internal/config"
	"golang/model"
	"golang/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.UserRepository
	tokenManager ITokenManager
}

func NewUserService(repo repository.UserRepository, tm ITokenManager) *UserService {
	return &UserService{repo: repo, tokenManager: tm}
}

func (s *UserService) SignUp(ctx context.Context, credential *model.Credential) error {
	if err := credential.Validate(); err != nil {
		return err
	}

	hashedPassword, err := config.HashPassword(credential.Password)

	if err != nil {
		return err
	}

	credential.Password = string(hashedPassword)

	return s.repo.SignUp(ctx, credential)
}

func (s *UserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)

	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))

	if err != nil {
		return nil, model.ErrInvalidCredentials
	}

	return user, nil
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	user, err := s.tokenManager.VerifyRefreshToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	newAccess, err := s.tokenManager.CreateAccessToken(&user)

	if err != nil {
		return "", "", err
	}

	newRefresh, err := s.tokenManager.CreateRefreshToken(user)

	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil

}
