package token

import (
	"context"
	"golang/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type MockUserRepository struct {
	GetByEmailFunc func(ctx context.Context, email string) (*model.User, error)
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := &MockUserRepository{}

	os.Setenv("JWT_SECRET_KEY", "test-secret")

	mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*model.User, error) {
		return &model.User{Email: "test@test.com", SecurityStamp: "some-uuid"}, nil
	}

	t.Run("Missing Authorization Header", func(t *testing.T) {

		r := gin.New()

		r.GET("/test", Middleware(mockRepo), func(c *gin.Context) {
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
		r.GET("/test", Middleware(mockRepo), func(c *gin.Context) {
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

	mockRepo := &MockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				Email:         "doe@doe.com",
				SecurityStamp: "test-stamp",
			}, nil
		},
	}

	claims := jwt.MapClaims{
		"user_name":      "John",
		"user_email":     "doe@doe.com",
		"security_stamp": "test-stamp",
		"exp":            time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte("test-secret"))

	r := gin.New()
	nextCalled := false

	r.GET("/test", Middleware(mockRepo), func(c *gin.Context) {
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

func TestRequireRole_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	calledNext := false

	r.POST("/tasks", func(c *gin.Context) {
		c.Set("user", model.User{Role: "ADMIN"})
		c.Next()
	}, RequireRole("ADMIN"), func(c *gin.Context) {
		calledNext = true
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
	if !calledNext {
		t.Error("Next handler was not called")
	}
}

func TestRequireRole_Manual(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		rec := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(rec)
		ctx.Set("user", model.User{Role: "ADMIN"})

		RequireRole("ADMIN")(ctx)

		if ctx.IsAborted() {
			t.Error("Expected middleware to proceed, but it aborted")
		}
	})
}
