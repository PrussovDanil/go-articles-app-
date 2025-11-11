package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)


func LoadConfig() (Config, error) {

	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("failed to load .env file: %w", err)
	}

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT value: %w", err)
	}

	cfg := Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.DBName == "" {
		return Config{}, fmt.Errorf("missing required database configuration")
	}

	return cfg, nil
}
