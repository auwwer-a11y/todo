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
	"time"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logger.New()

	pgConn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.TasksPostgres.User, cfg.TasksPostgres.Password, cfg.TasksPostgres.Host, cfg.TasksPostgres.Port, cfg.TasksPostgres.DBName)
	pgClient, err := sqlx.Connect("postgres", pgConn)
	if err != nil {
		panic(err)
	}
	taskRepo := postgres.NewTaskRepo(pgClient)

	mongoURI := fmt.Sprintf("mongodb://%s:%s",
    cfg.Mongo.Host, cfg.Mongo.Port)

	mongoClient, err := mongodriver.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	notificationRepo := mongoadapter.NewNotificationRepo(mongoClient.Database(cfg.Mongo.DBName))

	notifierService := service.NewNotifierService(notificationRepo, taskRepo, log)

	consumer := brokeradapter.NewKafkaConsumer(log, cfg.Kafka.Brokers, cfg.Kafka.Topic, "notifier-group")
	go func() {
    	for {
       	 msg, err := consumer.ReadMessage(context.Background())
        	if err != nil {
          	  log.Error("Failed to read message", "error", err)
           		continue
       	 }
			var event map[string]interface{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Error("Failed to unmarshal message", "error", err)
				continue
			}

			eventType, _ := event["event_type"].(string)
			log.Info("Received event", "event_type", eventType)
			if eventType == "task.status_changed" {
    			payload, _ := event["payload"].(map[string]interface{})
    			taskID, _ := payload["id"].(string)
    			userID, _ := payload["user_id"].(string)
    			status, _ := payload["status"].(string)
    			notifierService.HandleStatusChanged(context.Background(), taskID, userID, status)
}
    	}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := notifierService.CheckDeadlines(context.Background()); err != nil {
					log.Error("Failed to check deadlines", "error", err)
				}
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down notifier service")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoClient.Disconnect(ctx)
	pgClient.Close()
	consumer.Close()
}