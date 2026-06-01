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
		slog.Info("Decode payload error", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Info("Invalid Data", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	err := c.service.SignUp(ctx.Request.Context(), &credentials)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (c *UserController) Login(ctx *gin.Context) {
	var credentials model.Credential

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		slog.Info("Decode payload error", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Info("Invalid Data", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Data"})
		return
	}

	user, err := c.service.Login(ctx.Request.Context(), credentials.Email, credentials.Password)

	if err != nil {
		slog.Warn("Failed login attempt", "email", credentials.Email, "err", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(user)

	if err != nil {
		slog.Error("Token generation failed", "err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (c *UserController) Refresh(ctx *gin.Context) {
	var req model.Request

	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	newAccess, newRefresh, err := c.service.Refresh(ctx.Request.Context(), req.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or invalid token"})
		return
	}

	ctx.JSON(http.StatusOK, model.Response{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}
