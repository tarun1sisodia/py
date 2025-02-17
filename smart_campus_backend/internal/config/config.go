package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	JWT struct {
		Secret string
	}
	Firebase struct {
		CredentialsPath string
	}
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{}

	// Database configuration
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")

	// JWT configuration
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	// Firebase configuration
	cfg.Firebase.CredentialsPath = os.Getenv("FIREBASE_CREDENTIALS_PATH")

	return cfg, nil
}
