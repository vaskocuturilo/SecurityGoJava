package main

import (
	"context"
	"database/sql"
	"errors"
	"golang/app"
	"golang/controller"
	"golang/internal/config"
	"golang/migrations"
	"golang/repository"
	"golang/service"
	"golang/token"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.Postgres.ConnString())

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Warn("Failed to init Postgres: ", "error", err)
		}
	}(db)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.DBTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Warn("Postgres is not ready: ", "error", err)
	} else {
		slog.Info("Postgres is ready")
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		slog.Warn("Failed to init migration process: ", "error", err)
	} else {
		slog.Info("Successfully init migration process")
	}

	repo := repository.NewPostgresUserRepository(db)

	serv := service.NewUserService(repo)

	ctrl := controller.NewEventController(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/auth/register", ctrl.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", ctrl.Login)
	mux.HandleFunc("POST /api/v1/auth/refresh", ctrl.Refresh)
	mux.Handle("GET /api/v1/tasks", token.Middleware(http.HandlerFunc(app.Tasks)))

	srv := http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: mux}

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
