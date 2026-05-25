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
)

type MockUserService struct {
	SignUpFunc  func(ctx context.Context, credential *model.Credential) error
	LoginFunc   func(ctx context.Context, email, password string) (*model.User, error)
	RefreshFunc func(ctx context.Context, refreshToken string) (string, string, error)
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

func (m *MockUserService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	if m.RefreshFunc != nil {
		return m.RefreshFunc(ctx, refreshToken)
	}
	return "", "", nil
}

func TestUserController_SignUp_Table(t *testing.T) {
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

			var buf bytes.Buffer
			json.NewEncoder(&buf).Encode(tc.body)

			req := httptest.NewRequest(http.MethodPost, "/register", &buf)
			rec := httptest.NewRecorder()

			ctrl.SignUp(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}
		})
	}
}

func TestUserController_Login_Table(t *testing.T) {
	testUser := &model.User{
		ID:             "550e8400-e29b-41d4-a716-446655440100",
		Email:          "test@test.com",
		Username:       "hashed-password",
		HashedPassword: "Description",
	}

	tests := []struct {
		name         string
		giveEmail    string
		givePassword string
		wantReturn   *model.User
		mockErr      error
		wantStatus   int
	}{
		{
			name:         "Success",
			giveEmail:    testUser.Email,
			givePassword: testUser.HashedPassword,
			wantReturn:   testUser,
			mockErr:      nil,
			wantStatus:   http.StatusOK,
		},
		{
			name:         "Validation Incorrect Email",
			giveEmail:    "",
			givePassword: testUser.HashedPassword,
			wantReturn:   nil,
			mockErr:      model.ErrEmailRequired,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "Validation Not Found",
			giveEmail:    "test1@test.com",
			givePassword: testUser.HashedPassword,
			wantReturn:   nil,
			mockErr:      model.ErrUserNotFound,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "Validation Internal Server Error",
			giveEmail:    testUser.Email,
			givePassword: testUser.HashedPassword,
			wantReturn:   nil,
			mockErr:      errors.New("something went wrong in the database"),
			wantStatus:   http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &MockUserService{
				LoginFunc: func(ctx context.Context, email, password string) (*model.User, error) {
					return tc.wantReturn, tc.mockErr
				},
			}
			ctrl := NewUserController(mockService)

			req := httptest.NewRequest(http.MethodPost, "/login/", nil)

			if tc.giveEmail != "" || tc.givePassword != "" {
				req.SetBasicAuth(tc.giveEmail, tc.givePassword)
			}

			rec := httptest.NewRecorder()

			ctrl.Login(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tc.wantStatus)
			}

			if tc.wantStatus == http.StatusOK {
				var response map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("failed to unmarshal: %v", err)
				}
				if _, ok := response["access_token"]; !ok {
					t.Error("expected access_token in response body")
				}
			}
		})
	}
}
