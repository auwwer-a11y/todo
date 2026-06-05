package main

import (
	_ "github.com/lib/pq"
	"context"
	mongoadapter "github.com/auwwer-a11y/todo/internal/adapter/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"github.com/jmoiron/sqlx"
	"github.com/auwwer-a11y/todo/pkg/config"
	"github.com/auwwer-a11y/todo/pkg/logger"
	"github.com/auwwer-a11y/todo/internal/adapter/postgres"
	"github.com/auwwer-a11y/todo/internal/service"
	"fmt"
	brokeradapter "github.com/auwwer-a11y/todo/internal/adapter/broker"
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

	// Initalize MongoDB
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
		cfg.TasksPostgres.User, cfg.TasksPostgres.Password, cfg.TasksPostgres.Host, cfg.TasksPostgres.Port, cfg.TasksPostgres.DBName)
	pgClient, err := sqlx.Connect("postgres", pgConn)
	if err != nil {
		panic(err)
	}
	if err := goose.Up(pgClient.DB, "migrations/tasks"); err != nil {
		panic(err)
	}
	
	taskRepo := postgres.NewTaskRepo(pgClient)

	// initialize services
	kafkaProducer := brokeradapter.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	taskService := service.NewTaskService(taskRepo, noteRepo, kafkaProducer, log)
	noteService := service.NewNoteService(noteRepo, taskRepo, log)

	router := handler.NewTasksRouter(taskService, noteService, log, cfg.App.AuthServiceURL)

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

	mongoClient.Disconnect(context.Background())
	pgClient.Close()
}