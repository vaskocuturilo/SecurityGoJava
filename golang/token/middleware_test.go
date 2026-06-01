package token

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("JWT_SECRET_KEY", "test-secret")

	t.Run("Missing Authorization Header", func(t *testing.T) {
		r := gin.New()

		r.GET("/test", Middleware(), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})

	t.Run("Invalid Token Format", func(t *testing.T) {
		r := gin.New()
		r.GET("/test", Middleware(), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Basic some-token")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", rec.Code)
		}
	})
}

func TestMiddleware_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("JWT_SECRET_KEY", "test-secret")

	claims := jwt.MapClaims{
		"user_name":  "John",
		"user_email": "doe@doe.com",
		"exp":        time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	r := gin.New()
	nextCalled := false

	r.GET("/test", Middleware(), func(c *gin.Context) {
		nextCalled = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	if !nextCalled {
		t.Error("Next handler was not called")
	}
}
