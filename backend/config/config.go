package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OTP      OTPConfig
}

type ServerConfig struct {
	Port         string
	Environment  string
	AllowOrigins []string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret        string
	TokenExpiry   time.Duration
	RefreshSecret string
	RefreshExpiry time.Duration
}

type OTPConfig struct {
	Length      int
	Expiry      time.Duration
	MaxAttempts int
}

var AppConfig Config

func Load() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Server Configuration
	AppConfig.Server = ServerConfig{
		Port:         getEnv("PORT", "8080"),
		Environment:  getEnv("GIN_MODE", "debug"),
		AllowOrigins: []string{"*"}, // Update this for production
	}

	// Database Configuration
	AppConfig.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "smart_attendance"),
	}

	// JWT Configuration
	tokenExpiry, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		tokenExpiry = 24
	}
	refreshExpiry, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))
	if err != nil {
		refreshExpiry = 168 // 7 days
	}

	AppConfig.JWT = JWTConfig{
		Secret:        getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiry:   time.Duration(tokenExpiry) * time.Hour,
		RefreshSecret: getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
		RefreshExpiry: time.Duration(refreshExpiry) * time.Hour,
	}

	// OTP Configuration
	otpLength, err := strconv.Atoi(getEnv("OTP_LENGTH", "6"))
	if err != nil {
		otpLength = 6
	}
	otpExpiry, err := strconv.Atoi(getEnv("OTP_EXPIRY_MINUTES", "5"))
	if err != nil {
		otpExpiry = 5
	}
	maxAttempts, err := strconv.Atoi(getEnv("OTP_MAX_ATTEMPTS", "3"))
	if err != nil {
		maxAttempts = 3
	}

	AppConfig.OTP = OTPConfig{
		Length:      otpLength,
		Expiry:      time.Duration(otpExpiry) * time.Minute,
		MaxAttempts: maxAttempts,
	}

	return nil
}

func GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.Database.User,
		AppConfig.Database.Password,
		AppConfig.Database.Host,
		AppConfig.Database.Port,
		AppConfig.Database.Name,
	)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
