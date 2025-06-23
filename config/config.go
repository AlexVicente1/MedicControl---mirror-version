package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	JWTSecret  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	RateLimit  int
	LogDir     string
	StaticDir  string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		Port:       os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		RateLimit:  getEnvInt("RATE_LIMIT", 100),
		LogDir:     os.Getenv("LOG_DIR"),
		StaticDir:  os.Getenv("STATIC_DIR"),
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}
