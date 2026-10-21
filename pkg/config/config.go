package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL      string
	DBHost     string
	DBPort     string
	DBSSLMode  string
	DBName     string
	DBUser     string
	DBPassword string
}

func Load() *Config {
	// Локальный .env необязателен: на хостинге
	// переменные окружения задаются через панель.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: could not load .env: %v", err)
	}

	return &Config{
		DBURL:      os.Getenv("DB_URL"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}
