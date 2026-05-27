package token

import (
	"context"
	"errors"
	"fmt"
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

func NewTokenManager(repo repository.UserRepository) *TokenManager {
	return &TokenManager{
		repo:          repo,
		refreshTokens: make(map[string]struct{}),
	}
}

func (m *TokenManager) CreateAccessToken(u *model.User) (string, error) {
	secret := config.JWTSecret()
	issuer := config.GetIssuer()

	tokenID := uuid.New().String()

	var mapClaims = jwt.MapClaims{
		"iss":       issuer,
		"sub":       u.ID,
		"jti":       tokenID,
		"nbf":       time.Now().Unix(),
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(config.AccessTokenDuration()).Unix(),
		"user_name": u.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	signedToken, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedToken, err
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

	m.mx.Lock()
	m.refreshTokens[tokenID] = struct{}{}
	m.mx.Unlock()

	return signed, nil
}

func (m *TokenManager) VerifyRefreshToken(refreshToken string) (model.User, error) {
	c := &claims.UserClaims{}

	token, err := jwt.ParseWithClaims(refreshToken, c, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecret(), nil
	})

	if err != nil || !token.Valid || c.TokenType != "refresh" {
		return model.User{}, ErrInvalidToken
	}

	m.mx.Lock()
	if _, exists := m.refreshTokens[c.ID]; !exists {
		m.mx.Unlock()
		return model.User{}, ErrInvalidToken
	}
	delete(m.refreshTokens, c.ID)
	m.mx.Unlock()

	userPtr, err := m.repo.GetByEmail(context.Background(), c.Subject)
	if err != nil {
		return model.User{}, err
	}

	if userPtr == nil {
		return model.User{}, ErrInvalidToken
	}

	return *userPtr, nil
}
