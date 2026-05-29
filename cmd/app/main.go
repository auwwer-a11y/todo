package main

import (
	"context"
	mongoadapter "github.com/auwwer-a11y/todo/internal/adapter/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"github.com/redis/go-redis/v9"
	"github.com/jmoiron/sqlx"
	"github.com/auwwer-a11y/todo/internal/config"
	"github.com/auwwer-a11y/todo/internal/adapter/postgres"
	"github.com/auwwer-a11y/todo/internal/logger"
	"github.com/auwwer-a11y/todo/internal/service"
	"fmt"
	redisadapter "github.com/auwwer-a11y/todo/internal/adapter/redis"
	brokeradapter "github.com/auwwer-a11y/todo/internal/adapter/broker"
	"net/http"
	"time"
	"github.com/auwwer-a11y/todo/internal/handler"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New()


	// Initialize MongoDB
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		cfg.Mongo.User, cfg.Mongo.Password, cfg.Mongo.Host, cfg.Mongo.Port)

	mongoClient, err := mongodriver.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	db := mongoClient.Database("todo")
	noteRepo := mongoadapter.NewNoteRepo(db)

	// Initialize PostgreSQL
	pgConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName)
	pgClient, err := sqlx.Connect("postgres", pgConn)
	if err != nil {
		panic(err)
	}
	userRepo := postgres.NewUserRepo(pgClient)
	taskRepo := postgres.NewTaskRepo(pgClient)

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
	})
	cacheRepo := redisadapter.NewCache(redisClient)

	// Initialize services
	kafkaProducer := brokeradapter.NewKafkaProducer([]string{cfg.Kafka.Brokers}, cfg.Kafka.Topic)
	ttl, _ := time.ParseDuration(cfg.App.JWTTTL)
	userService := service.NewUserService(userRepo, cacheRepo, log, cfg.App.JWTSecret, ttl)
	taskService := service.NewTaskService(taskRepo, noteRepo, kafkaProducer, log)
	noteService := service.NewNoteService(noteRepo, taskRepo, log)

	// Initialize router
	router := handler.NewRouter(userService, taskService, noteService, log, cfg.App.JWTSecret)

	// Start server
	srv := &http.Server{
		Addr: ":" + cfg.App.Port,
		Handler: router,
	}

	log.Info("Starting server", "port", cfg.App.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Server error", "error", err)
	}
}