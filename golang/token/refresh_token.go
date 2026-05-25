package token

import (
	"context"
	"errors"
	"golang/claims"
	"golang/internal/config"
	"golang/model"
	"golang/repository"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

var (
	mx sync.RWMutex

	refreshTokens = make(map[string]struct{})
)

type TokenManager struct {
	repo          repository.UserRepository
	mx            sync.RWMutex
	refreshTokens map[string]struct{}
}

func (m *TokenManager) CreateAccessToken(user *model.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewTokenManager(repo repository.UserRepository) *TokenManager {
	return &TokenManager{
		repo:          repo,
		refreshTokens: make(map[string]struct{}),
	}
}

func (m *TokenManager) CreateRefreshToken(u model.User) (string, error) {
	tokenID := uuid.New().String()

	userClaims := claims.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GetIssuer(),
			Subject:   u.Email,
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenType: "refresh",
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims).
		SignedString(config.JWTSecret())

	if err != nil {
		return "", err
	}

	mx.Lock()
	refreshTokens[tokenID] = struct{}{}
	mx.Unlock()

	return signed, nil
}

func (m *TokenManager) VerifyRefreshToken(refreshToken string) (model.User, error) {
	claims := &claims.UserClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecret(), nil
	})

	if err != nil || !token.Valid || claims.TokenType != "refresh" {
		return model.User{}, ErrInvalidToken
	}

	m.mx.Lock()
	defer m.mx.Unlock()

	if _, exists := m.refreshTokens[claims.ID]; !exists {
		return model.User{}, ErrInvalidToken
	}
	delete(m.refreshTokens, claims.ID)

	userPtr, err := m.repo.GetByEmail(context.Background(), claims.Email)

	if err != nil {
		return model.User{}, err
	}

	if userPtr == nil {
		return model.User{}, ErrInvalidToken
	}

	return *userPtr, nil

}
