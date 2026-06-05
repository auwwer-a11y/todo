package config

import (
	"os"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type AppConfig struct {
	Port string
	JWTSecret string
	JWTTL string
	AuthServiceURL string
}

type PostgresConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

type MongoConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

type RedisConfig struct {
	Host string
	Port string
	Password string
	DB int
}

type KafkaConfig struct {
	Brokers []string
	Topic string
}

type Config struct {
	App AppConfig
	Mongo MongoConfig
	Redis RedisConfig
	Kafka KafkaConfig
	AuthPostgres AuthPostgresConfig
	TasksPostgres TasksPostgresConfig
}

type AuthPostgresConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

type TasksPostgresConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		App: AppConfig{
			Port: os.Getenv("APP_PORT"),
			JWTSecret: os.Getenv("JWT_SECRET"),
			JWTTL: os.Getenv("JWT_TTL"),
			AuthServiceURL: os.Getenv("AUTH_SERVICE_URL"),
		},

		AuthPostgres: AuthPostgresConfig{
			Host: os.Getenv("AUTH_POSTGRES_HOST"),
			Port: os.Getenv("AUTH_POSTGRES_PORT"),
			User: os.Getenv("AUTH_POSTGRES_USER"),
			Password: os.Getenv("AUTH_POSTGRES_PASSWORD"),
			DBName: os.Getenv("AUTH_POSTGRES_DB"),
		},

		TasksPostgres: TasksPostgresConfig{
			Host: os.Getenv("TASKS_POSTGRES_HOST"),
			Port: os.Getenv("TASKS_POSTGRES_PORT"),
			User: os.Getenv("TASKS_POSTGRES_USER"),
			Password: os.Getenv("TASKS_POSTGRES_PASSWORD"),
			DBName: os.Getenv("TASKS_POSTGRES_DB"),
		},

		Mongo: MongoConfig{
			Host: os.Getenv("MONGO_HOST"),
			Port: os.Getenv("MONGO_PORT"),
			User: os.Getenv("MONGO_USER"),
			Password	: os.Getenv("MONGO_PASSWORD"),
			DBName: os.Getenv("MONGO_DB"),
		},
		Redis: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB: 0, // Default Redis DB
		},
		Kafka: KafkaConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKERS")},
			Topic: os.Getenv("KAFKA_TOPIC"),
		},
	}

	return config, nil
}