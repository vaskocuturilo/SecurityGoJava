package controller

import (
	"encoding/json"
	"errors"
	"golang/model"
	"golang/service"
	"golang/token"
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
	email, password, ok := r.BasicAuth()

	if !ok {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.service.Login(r.Context(), email, password)

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, err := token.CreateAccessToken(*user)
	if err != nil {
		slog.Error("JWT creation failed", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func (c *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	//TODO: refresh functionality
}
