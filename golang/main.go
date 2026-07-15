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

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

func main() {
	limiter := tollbooth.NewLimiter(5, nil)
	limiter.SetMessage("Too many requests, try again later.")

	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.Postgres.ConnString())

	if err != nil {
		slog.Error("Failed to init Postgres: ", "error", err)
		os.Exit(1)
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.DBTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		slog.Error("Postgres is not ready: ", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Postgres is ready")
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		slog.Error("Failed to init migration process: ", "error", err)
		os.Exit(1)
	} else {
		slog.Info("Successfully init migration process")
	}

	repo := repository.NewPostgresUserRepository(db)

	tokenMgr := token.NewTokenManager(repo)

	serv := service.NewUserService(repo, tokenMgr)

	ctrl := controller.NewUserController(serv)

	r := gin.Default()

	api := r.Group("/api/v1")
	users := api.Group("/users")

	users.POST("/register", tollbooth_gin.LimitHandler(limiter), ctrl.SignUp)
	users.POST("/login", tollbooth_gin.LimitHandler(limiter), ctrl.Login)
	users.POST("/refresh", ctrl.Refresh)
	users.POST("/logout", ctrl.Logout)
	api.GET("/tasks", token.Middleware(repo), app.Tasks)
	api.POST("/tasks", token.Middleware(repo), token.RequireRole("CREATE"), app.CreateTask)

	srv := &http.Server{Addr: net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), Handler: r}

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

	ctx, cancel = context.WithTimeout(context.Background(), cfg.Server.TTL)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Warn("Server forced to shutdown: ", "error", err)
	}

	slog.Info("Server exited properly")
}
