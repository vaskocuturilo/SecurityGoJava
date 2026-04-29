package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"golang/internal/config"
	"golang/users"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var dummyHash []byte

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", register)
	mux.HandleFunc("POST /auth/login", login)
	mux.HandleFunc("POST /auth/refresh", refresh)
	mux.HandleFunc("GET /tasks", tasks)

	srv := http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port)}

	go func() {
		slog.Info("Server is starting on port", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Warn("Listen error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	slog.Info("Shutdown signal received, shutting down gracefully...")
}

func login(w http.ResponseWriter, r *http.Request) {
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

	user, err := AuthUser(creds[0], creds[1])

	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, err := createAccessToken(user)
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

func AuthUser(email, password string) (users.User, error) {
	var ErrInvalidUserOrPassword = errors.New("invalid credentials")

	userDB := config.UserDB()

	idx := slices.IndexFunc(userDB, func(u users.User) bool {
		return strings.EqualFold(email, u.Email)
	})

	var hashToCheck []byte
	var user users.User

	if idx == -1 {
		hashToCheck = dummyHash
	} else {
		user = userDB[idx]
		hashToCheck = user.HashedPassword
	}

	err := checkPassword(hashToCheck, password)

	if idx == -1 || err != nil {
		return users.User{}, ErrInvalidUserOrPassword
	}

	return user, nil
}

func createAccessToken(u users.User) (string, error) {
	secret := config.JWTSecret()
	issuer := config.GetIssuer()

	tokenID := uuid.New().String()

	var claims = jwt.MapClaims{
		"iss":           issuer,
		"sub":           u.ID,
		"jti":           tokenID,
		"nbf":           time.Now().Unix(),
		"iat":           time.Now().Unix(),
		"exp":           time.Now().Add(config.AccessTokenDuration()).Unix(),
		"user_name":     u.Name,
		"user_lastname": u.Lastname,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return signedToken, err
}

func checkPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func register(w http.ResponseWriter, r *http.Request) {
}

func refresh(w http.ResponseWriter, r *http.Request) {
}

func tasks(w http.ResponseWriter, r *http.Request) {
}
