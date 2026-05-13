package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
	JWT      JWTConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Host           string
	Port           string
	TTL            time.Duration
	RequestTimeout time.Duration
	RedisTimeout   time.Duration
	DBTimeout      time.Duration
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   getEnv("POSTGRES_DB", "tasks_db"),
		},

		Server: ServerConfig{
			Host:           getEnv("SERVER_HOST", "0.0.0.0"),
			Port:           getEnv("SERVER_PORT", "8080"),
			TTL:            time.Minute * 10,
			DBTimeout:      time.Second * 5,
			RedisTimeout:   time.Second * 5,
			RequestTimeout: time.Second * 5,
		},
		JWT: JWTConfig{
			Secret: os.Getenv("SECRET_KEY"),
		},
	}
}

func (p PostgresConfig) ConnString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
