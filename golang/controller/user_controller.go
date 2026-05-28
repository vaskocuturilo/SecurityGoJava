package controller

import (
	"encoding/json"
	"errors"
	"golang/model"
	"golang/service"
	"log/slog"
	"net/http"
)

type UserController struct {
	service service.IUserService
}

func NewUserController(service service.IUserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var credentials model.Credential

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		slog.Info("Decode payload error", "error", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Info("Invalid Data", "error", err)
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	err := c.service.SignUp(r.Context(), &credentials)

	if err != nil {
		if errors.Is(err, model.ErrAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var credentials model.Credential

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		slog.Info("Decode payload error", "error", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := credentials.Validate(); err != nil {
		slog.Info("Invalid Data", "error", err)
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}

	user, err := c.service.Login(r.Context(), credentials.Email, credentials.Password)

	if err != nil {
		slog.Warn("Failed login attempt", "email", credentials.Email, "err", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := c.service.GenerateTokens(user)

	if err != nil {
		slog.Error("Token generation failed", "err", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(model.Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (c *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req model.Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	newAccess, newRefresh, err := c.service.Refresh(r.Context(), req.RefreshToken)

	if err != nil {
		http.Error(w, "Unauthorized or invalid token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(model.Response{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}
