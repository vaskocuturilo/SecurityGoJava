package main

import (
	"context"
	"database/sql"
	"errors"
	"golang/auth"
	"golang/internal/config"
	"golang/migrations"
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

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth/register", auth.Register)
	mux.HandleFunc("POST /auth/login", auth.Login)
	mux.HandleFunc("POST /auth/refresh", auth.Refresh)
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

func tasks(w http.ResponseWriter, r *http.Request) {
	//TODO: add tasks functionality.
}
