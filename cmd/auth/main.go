package main

import (
	_ "github.com/lib/pq"
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/auwwer-a11y/todo/pkg/config"
	"github.com/auwwer-a11y/todo/pkg/logger"
	"github.com/auwwer-a11y/todo/internal/adapter/postgres"
	"github.com/auwwer-a11y/todo/internal/service"
	"fmt"
	nethttp "net/http"
	"time"
	"github.com/auwwer-a11y/todo/internal/handler"
	"os"
	"os/signal"
	"syscall"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New()


	pgConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.AuthPostgres.User, cfg.AuthPostgres.Password, cfg.AuthPostgres.Host, cfg.AuthPostgres.Port, cfg.AuthPostgres.DBName)
	pgClient, err := sqlx.Connect("postgres", pgConn)
	if err != nil {
		panic(err)
	}
	if err := goose.Up(pgClient.DB, "migrations/auth"); err != nil {
		panic(err)
	}
	userRepo := postgres.NewUserRepo(pgClient)

	ttl, _ := time.ParseDuration(cfg.App.JWTTL)
	userService := service.NewUserService(userRepo, nil, log, cfg.App.JWTSecret, ttl)

	router := handler.NewAuthRouter(userService, log, cfg.App.JWTSecret)

	srv := &nethttp.Server{
		Addr: ":" + cfg.App.Port,
		Handler: router,
	}

	go func() {
		log.Info("Starting server", "port", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != nethttp.ErrServerClosed {
			log.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	pgClient.Close()
}

