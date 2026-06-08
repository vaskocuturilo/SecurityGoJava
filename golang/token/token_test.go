package token

import (
	"context"
	"golang/claims"
	"golang/internal/config"
	"golang/model"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func (m *MockUserRepository) Logout(ctx context.Context, userID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) SignUp(ctx context.Context, credential *model.Credential) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) Refresh(refreshToken string) (string, string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.GetByEmailFunc(ctx, email)
}

func TestTokenManager_CreateAccessToken(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "super-secret-test-key-123")

	tokenMgr := NewTokenManager(nil)

	mockUser := model.User{
		ID:       "f47ac10b-58cc-4372-a567",
		Username: "John",
		Email:    "Doe@doe.com",
		Role:     "USER",
	}

	tokenString, err := tokenMgr.CreateAccessToken(&mockUser)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.JWTSecret(), nil
	})

	if err != nil || !token.Valid {
		t.Errorf("Token is not valid: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Could not parse claims")
	}

	if claims["sub"] != mockUser.ID {
		t.Errorf("Expected sub %s, got %s", mockUser.ID, claims["sub"])
	}
}

func TestTokenManager_CreateRefreshTokenToken(t *testing.T) {
	tokenMgr := NewTokenManager(nil)

	mockUser := model.User{
		Email: "Doe@doe.com",
	}

	tokenString, err := tokenMgr.CreateRefreshToken(mockUser)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	c := &claims.UserClaims{}
	_, err = jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (interface{}, error) {
		return config.JWTSecret(), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse created token: %v", err)
	}

	tokenMgr.mx.RLock()
	_, exists := tokenMgr.refreshTokens[c.ID]
	tokenMgr.mx.RUnlock()

	if !exists {
		t.Error("Expected refresh token to be stored in the manager's map")
	}

	if c.TokenType != "refresh" {
		t.Errorf("Expected token type 'refresh', got %s", c.TokenType)
	}
}

func TestTokenManager_VerifyRefreshToken(t *testing.T) {
	mockRepo := &MockUserRepository{}
	tokenMgr := NewTokenManager(mockRepo)

	mockUser := model.User{Email: "doe@doe.com"}

	mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return &mockUser, nil
	}

	tokenString, err := tokenMgr.CreateRefreshToken(mockUser)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	user, err := tokenMgr.VerifyRefreshToken(tokenString)
	if err != nil {
		t.Fatalf("Expected successful verification, got error: %v", err)
	}
	if user.Email != mockUser.Email {
		t.Errorf("Expected email %s, got %s", mockUser.Email, user.Email)
	}

	_, err = tokenMgr.VerifyRefreshToken(tokenString)
	if err != ErrInvalidToken {
		t.Errorf("Expected ErrInvalidToken on second use, got: %v", err)
	}
}
