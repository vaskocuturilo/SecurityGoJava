package service

import (
	"context"
	"errors"
	"golang/model"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	SignUpFunc     func(ctx context.Context, credential *model.Credential) error
	GetByEmailFunc func(ctx context.Context, email string) (*model.User, error)
	RefreshFunc    func(refreshToken string) (string, string, error)
}

type MockTokenManager struct {
	CreateAccessTokenFunc  func(user *model.User) (string, error)
	CreateRefreshTokenFunc func(u model.User) (string, error)
	VerifyRefreshTokenFunc func(refreshToken string) (model.User, error)
}

func (m *MockRepository) SignUp(ctx context.Context, credential *model.Credential) error {
	if m.SignUpFunc != nil {
		return m.SignUpFunc(ctx, credential)
	}
	return nil
}

func (m *MockRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockRepository) Refresh(refreshToken string) (string, string, error) {
	if m.RefreshFunc != nil {
		return m.RefreshFunc(refreshToken)
	}
	return "", "", nil
}

func (m *MockTokenManager) CreateAccessToken(u *model.User) (string, error) {
	return m.CreateAccessTokenFunc(u)
}
func (m *MockTokenManager) CreateRefreshToken(u model.User) (string, error) {
	return m.CreateRefreshTokenFunc(u)
}
func (m *MockTokenManager) VerifyRefreshToken(token string) (model.User, error) {
	return m.VerifyRefreshTokenFunc(token)
}

func TestUserService_SignUp_TableDriven(t *testing.T) {
	mockTM := &MockTokenManager{}

	type testCase struct {
		name         string
		giveEmail    string
		givePassword string
		mockResponse error
		wantErr      error
	}

	tests := []testCase{
		{
			name:         "Success",
			giveEmail:    "test@test.com",
			givePassword: "hashed-password",
			mockResponse: nil,
			wantErr:      nil,
		},
		{
			name:         "Empty Email - Validation Error",
			giveEmail:    "",
			givePassword: "hashed-password",
			mockResponse: nil,
			wantErr:      model.ErrEmailRequired,
		},
		{
			name:         "Empty Password - Validation Error",
			giveEmail:    "test1@test.com",
			givePassword: "",
			mockResponse: nil,
			wantErr:      model.ErrPasswordRequired,
		},
		{
			name:         "Password Too Short - Validation Error",
			giveEmail:    "test1@test.com",
			givePassword: "short",
			mockResponse: nil,
			wantErr:      model.ErrPasswordTooShort,
		},
		{
			name:         "Repository Conflict",
			giveEmail:    "Duplicate",
			givePassword: "hashed-password",
			mockResponse: model.ErrAlreadyExists,
			wantErr:      model.ErrAlreadyExists,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &MockRepository{
				SignUpFunc: func(ctx context.Context, credential *model.Credential) error {
					return tc.mockResponse
				},
			}

			serv := NewUserService(mockRepo, mockTM)
			user := &model.Credential{Email: tc.giveEmail, Password: tc.givePassword}

			// Act
			err := serv.SignUp(context.Background(), user)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}
		})
	}
}

func TestUserService_Login_TableDriven(t *testing.T) {
	mockTM := &MockTokenManager{}
	correctPassword := "super-secret-password"
	wrongPassword := "wrong-password"

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash for test: %v", err)
	}
	correctHash := string(hashedBytes)

	dbUser := model.User{
		ID:             "123",
		Email:          "test@test.com",
		HashedPassword: correctHash,
	}

	type testCase struct {
		name         string
		giveEmail    string
		givePassword string
		mockReturn   *model.User
		mockErr      error
		wantUser     *model.User
		wantErr      error
	}

	tests := []testCase{
		{
			name:         "Success",
			giveEmail:    "test@test.com",
			givePassword: correctPassword,
			mockReturn:   &dbUser,
			mockErr:      nil,
			wantUser:     &dbUser,
			wantErr:      nil,
		},
		{
			name:         "User Not Found In DB",
			giveEmail:    "notfound@test.com",
			givePassword: correctPassword,
			mockReturn:   nil,
			mockErr:      model.ErrUserNotFound,
			wantUser:     nil,
			wantErr:      model.ErrInvalidCredentials,
		},
		{
			name:         "Wrong Password",
			giveEmail:    "test@test.com",
			givePassword: wrongPassword,
			mockReturn:   &dbUser,
			mockErr:      nil,
			wantUser:     nil,
			wantErr:      model.ErrInvalidCredentials,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			mockRepo := &MockRepository{
				GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
					if email != tc.giveEmail {
						t.Errorf("Mock expected Email %s, got %s", tc.giveEmail, email)
					}
					return tc.mockReturn, tc.mockErr
				},
			}

			serv := NewUserService(mockRepo, mockTM)

			// Act
			result, err := serv.Login(context.Background(), tc.giveEmail, tc.givePassword)

			// Assert
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Expected error %v, but got %v", tc.wantErr, err)
			}

			if tc.wantUser != nil {
				if result == nil {
					t.Fatal("Expected result user, but got nil")
				}
				if result.ID != tc.wantUser.ID {
					t.Errorf("Expected user ID %s, got %s", tc.wantUser.ID, result.ID)
				}
				if result.Email != tc.wantUser.Email {
					t.Errorf("Expected user Email %s, got %s", tc.wantUser.Email, result.Email)
				}
				if result.HashedPassword != tc.wantUser.HashedPassword {
					t.Errorf("Expected user HashedPassword %s, got %s", tc.wantUser.HashedPassword, result.HashedPassword)
				}
			} else {
				if result != nil {
					t.Errorf("Expected nil user result, but got %v", result)
				}
			}
		})
	}
}

func TestUserService_GenerateTokens(t *testing.T) {
	mockTM := &MockTokenManager{
		CreateAccessTokenFunc: func(user *model.User) (string, error) {
			return "access-token", nil
		},
		CreateRefreshTokenFunc: func(u model.User) (string, error) {
			return "refresh-token", nil
		},
	}

	serv := NewUserService(nil, mockTM)
	user := &model.User{Email: "test@test.com"}

	access, refresh, err := serv.GenerateTokens(user)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if access != "access-token" || refresh != "refresh-token" {
		t.Errorf("Tokens mismatch: got %s, %s", access, refresh)
	}
}

func TestUserService_Refresh(t *testing.T) {
	mockUser := model.User{Email: "test@test.com"}
	mockTM := &MockTokenManager{
		VerifyRefreshTokenFunc: func(token string) (model.User, error) {
			if token != "valid-token" {
				return model.User{}, errors.New("invalid")
			}
			return mockUser, nil
		},
		CreateAccessTokenFunc: func(user *model.User) (string, error) {
			return "new-access", nil
		},
		CreateRefreshTokenFunc: func(u model.User) (string, error) {
			return "new-refresh", nil
		},
	}

	serv := NewUserService(nil, mockTM)

	t.Run("Success", func(t *testing.T) {
		access, refresh, err := serv.Refresh("valid-token")

		if err != nil {
			t.Errorf("Expected success, got %v", err)
		}
		if access != "new-access" || refresh != "new-refresh" {
			t.Errorf("Tokens mismatch")
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		_, _, err := serv.Refresh("bad-token")
		if err == nil {
			t.Error("Expected error for bad token, got nil")
		}
	})
}
