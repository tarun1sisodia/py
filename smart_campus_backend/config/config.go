package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Firebase FirebaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port           string
	Mode           string
	AllowedOrigins []string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type FirebaseConfig struct {
	CredentialsPath string
	ProjectID       string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
	RefreshSecret   string
	RefreshExpHours int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port:           getEnv("PORT", "8080"),
			Mode:           getEnv("GIN_MODE", "debug"),
			AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
			ReadTimeout:    time.Second * 10,
			WriteTimeout:   time.Second * 10,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 3306),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root1234"),
			DBName:   getEnv("DB_NAME", "smart_campus"),
		},
		Firebase: FirebaseConfig{
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", "config/firebase-credentials.json"),
			ProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			RefreshExpHours: getEnvAsInt("JWT_REFRESH_EXPIRATION_HOURS", 168),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
