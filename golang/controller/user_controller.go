package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang/model"
	"golang/service"
	"golang/token"
	"log/slog"
	"net/http"
	"strings"
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
	authHeader := r.Header.Get("Authorization")

	const basicAuthPrefix = "Basic "

	if authHeader == "" || !strings.HasPrefix(authHeader, basicAuthPrefix) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	payload, err := base64.StdEncoding.DecodeString(authHeader[len(basicAuthPrefix):])

	if err != nil {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
		return
	}

	creds := strings.SplitN(string(payload), ":", 2)

	if len(creds) != 2 {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
		return
	}

	user, err := token.AuthUser(creds[0], creds[1])

	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := token.CreateAccessToken(user)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	type response struct {
		AccessToken string `json:"access_token"`
	}

	err = json.NewEncoder(w).Encode(response{AccessToken: accessToken})
	if err != nil {
		return
	}
}

func (c *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	//TODO: refresh functionality
}
