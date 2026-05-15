package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang/model"
	"golang/service"
	"golang/token"
	"log/slog"
	"net/http"
	"strings"
)

const (
	userNotFound        = "User not found"
	notFound            = "Not found with ID"
	unexpected          = "Unexpected error"
	invalidUrl          = "Invalid ID in URL"
	internalServerError = "Internal server error"
)

type UserController struct {
	service service.IUserService
}

func NewEventController(service service.IUserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	var credential model.Credential

	if err := json.NewDecoder(r.Body).Decode(&credential); err != nil {
		slog.Info("Decode payload error", "error", err)
		http.Error(w, "Failed to Decode payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	newEvent := model.NewUser(credential.UserName, credential.UserName, credential.Password)

	err := c.service.Login(ctx, newEvent)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidInput):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, model.ErrAlreadyExists):
			http.Error(w, "Event ID already taken", http.StatusConflict)
		default:
			slog.Error(unexpected, "error", err)
			http.Error(w, internalServerError, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/events/%s", newEvent.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
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
