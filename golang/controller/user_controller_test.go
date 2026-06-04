package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"golang/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct {
	SignUpFunc         func(ctx context.Context, credential *model.Credential) error
	LoginFunc          func(ctx context.Context, email, password string) (*model.User, error)
	RefreshFunc        func(refreshToken string) (string, string, error)
	GenerateTokensFunc func(user *model.User) (string, string, error)
}

func (m *MockUserService) SignUp(ctx context.Context, credential *model.Credential) error {
	if m.SignUpFunc != nil {
		return m.SignUpFunc(ctx, credential)
	}
	return nil
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (*model.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, email, password)
	}
	return nil, nil
}

func (m *MockUserService) Refresh(refreshToken string) (string, string, error) {
	if m.RefreshFunc != nil {
		return m.RefreshFunc(refreshToken)
	}
	return "", "", nil
}

func (m *MockUserService) GenerateTokens(user *model.User) (string, string, error) {
	if m.GenerateTokensFunc != nil {
		return m.GenerateTokensFunc(user)
	}
	return "", "", nil
}

func TestUserController_SignUp_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       any
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			body:       model.Credential{Email: "test@test.com", Password: "hashed-password"},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Validation Already Exists",
			body:       model.Credential{Email: "test@test.com", Password: "hashed-password"},
			mockErr:    model.ErrAlreadyExists,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Validation Email Error",
			body:       model.Credential{Email: "", Password: "hashed-password"},
			mockErr:    model.ErrEmailRequired,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Password Error",
			body:       model.Credential{Email: "test@test.com", Password: ""},
			mockErr:    model.ErrPasswordRequired,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Password Too Short",
			body:       model.Credential{Email: "test@test.com", Password: "short"},
			mockErr:    model.ErrPasswordTooShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Invalid JSON",
			body:       "invalid",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Internal Server Error",
			body:       model.Credential{Email: "test@test.com", Password: "hashed-password"},
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockService := &MockUserService{
				SignUpFunc: func(ctx context.Context, credential *model.Credential) error { return tc.mockErr },
			}
			ctrl := NewUserController(mockService)

			rec := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rec)

			body, _ := json.Marshal(tc.body)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			ctrl.SignUp(ctx)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestUserController_Login_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUser := &model.User{
		ID:             "550e8400-e29b-41d4-a716-446655440100",
		Email:          "test@test.com",
		Username:       "hashed-password",
		HashedPassword: "Description",
	}

	tests := []struct {
		name       string
		body       any
		wantReturn *model.User
		mockErr    error
		wantStatus int
	}{
		{
			name:       "Success",
			body:       model.Credential{Email: "test@test.com", Password: "hashed-password"},
			wantReturn: testUser,
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Validation Incorrect Email",
			body:       model.Credential{Email: "", Password: "hashed-password"},
			wantReturn: nil,
			mockErr:    model.ErrEmailRequired,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Validation Not Found",
			body:       model.Credential{Email: "test1@test.com", Password: "hashed-password"},
			wantReturn: nil,
			mockErr:    model.ErrUserNotFound,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Validation Internal Server Error",
			body:       model.Credential{Email: "test@test.com", Password: "hashed-password"},
			wantReturn: nil,
			mockErr:    errors.New("something went wrong in the database"),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Validation Invalid JSON format",
			body:       "not-a-json",
			wantReturn: nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockUserService{
				LoginFunc: func(ctx context.Context, email, password string) (*model.User, error) {
					return tc.wantReturn, tc.mockErr
				},
				GenerateTokensFunc: func(user *model.User) (string, string, error) {
					return "access-token-val", "refresh-token-val", nil
				},
			}

			ctrl := NewUserController(mockService)

			rec := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(rec)

			body, _ := json.Marshal(tc.body)

			ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			ctx.Request.Header.Set("Content-Type", "application/json")

			ctrl.Login(ctx)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantStatus == http.StatusOK {
				var response map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if response["access_token"] != "access-token-val" {
					t.Error("expected correct access_token")
				}
			}
		})
	}
}

func TestUserController_Refresh(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := &MockUserService{
		RefreshFunc: func(refreshToken string) (string, string, error) {
			if refreshToken == "valid" {
				return "new-access", "new-refresh", nil
			}
			return "", "", errors.New("invalid")
		},
	}
	ctrl := NewUserController(mockService)

	t.Run("Success", func(t *testing.T) {
		rec := httptest.NewRecorder()

		ctx, _ := gin.CreateTestContext(rec)

		body := map[string]string{"refresh_token": "valid"}
		jsonBody, _ := json.Marshal(body)

		ctx.Request = httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(jsonBody))
		ctx.Request.Header.Set("Content-Type", "application/json")

		ctrl.Refresh(ctx)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rec.Code)
		}
	})
}
