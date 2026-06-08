package controller

import (
	"errors"
	"golang/model"
	"golang/service"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.IUserService
}

func NewUserController(service service.IUserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SignUp(ctx *gin.Context) {
	var credentials model.Credential

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		slog.Warn("Decode payload error", "error", err)
		errorResponse(ctx, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Warn("Validate data error", "error", err)
		errorResponse(ctx, http.StatusBadRequest, "Invalid Data")
		return
	}

	err := c.service.SignUp(ctx.Request.Context(), &credentials)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			errorResponse(ctx, http.StatusConflict, "User already exists")
			return
		}
		errorResponse(ctx, http.StatusInternalServerError, "Internal error")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (c *UserController) Login(ctx *gin.Context) {
	var credentials model.Credential

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		slog.Warn("Decode payload error", "error", err)
		errorResponse(ctx, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Warn("Invalid Data", "error", err)
		errorResponse(ctx, http.StatusBadRequest, "Invalid Data")
		return
	}

	user, err := c.service.Login(ctx.Request.Context(), credentials.Email, credentials.Password)

	if err != nil {
		slog.Warn("Failed login attempt", "email", credentials.Email, "err", err)
		errorResponse(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	slog.Info("User loaded from DB", "role", user.Role)

	accessToken, refreshToken, err := c.service.GenerateTokens(user)

	if err != nil {
		slog.Warn("Token generation failed", "err", err)
		errorResponse(ctx, http.StatusInternalServerError, "Internal error")
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (c *UserController) Refresh(ctx *gin.Context) {
	var req model.Request

	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, http.StatusBadRequest, "Invalid body")
		return
	}

	newAccess, newRefresh, err := c.service.Refresh(req.RefreshToken)

	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "Unauthorized or invalid token")
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}

func (c *UserController) Logout(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(string)
	err := c.service.Logout(ctx, userID)

	if err != nil {
		errorResponse(ctx, http.StatusInternalServerError, "Failed logout")
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out from all devices"})
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
