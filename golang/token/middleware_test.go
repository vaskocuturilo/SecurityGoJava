package token

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		Middleware(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Basic some-token") // Неверный формат (не Bearer)
		rec := httptest.NewRecorder()

		Middleware(nextHandler).ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})
}

func TestMiddleware_Success(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")

	claims := jwt.MapClaims{
		"user_name":  "John",
		"user_email": "doe@doe.com",
		"exp":        time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	nextCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	Middleware(nextHandler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	if !nextCalled {
		t.Error("Next handler was not called")
	}
}
