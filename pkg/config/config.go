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
	Postgres PostgresConfig
	Mongo MongoConfig
	Redis RedisConfig
	Kafka KafkaConfig
}

func LoadConfig() (*Config, error) {
	config := &Config{
		App: AppConfig{
			Port: os.Getenv("APP_PORT"),
			JWTSecret: os.Getenv("JWT_SECRET"),
			JWTTL: os.Getenv("JWT_TTL"),
		},
		Postgres: PostgresConfig{
			Host: os.Getenv("POSTGRES_HOST"),
			Port: os.Getenv("POSTGRES_PORT"),
			User: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName: os.Getenv("POSTGRES_DB"),
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