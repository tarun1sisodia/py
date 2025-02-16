package config

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

// Config represents the application configuration
type Config struct {
	Server     ServerConfig     `json:"server"`
	Database   DatabaseConfig   `json:"database"`
	JWT        JWTConfig        `json:"jwt"`
	Redis      RedisConfig      `json:"redis"`
	Validation ValidationConfig `json:"validation"`
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port            string          `json:"port"`
	Mode            string          `json:"mode"`
	Timeout         time.Duration   `json:"timeout"`
	ShutdownTimeout time.Duration   `json:"shutdown_timeout"`
	AllowedOrigins  []string        `json:"allowed_origins"`
	RateLimit       RateLimitConfig `json:"rate_limit"`
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// JWTConfig holds the JWT configuration
type JWTConfig struct {
	Secret           string        `json:"secret"`
	AccessTokenExp   time.Duration `json:"access_token_exp"`
	RefreshTokenExp  time.Duration `json:"refresh_token_exp"`
	SigningAlgorithm string        `json:"signing_algorithm"`
}

// RedisConfig holds the Redis configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

// ValidationConfig holds the validation rules
type ValidationConfig struct {
	PasswordMinLength int           `json:"password_min_length"`
	OTPLength         int           `json:"otp_length"`
	OTPExpiration     time.Duration `json:"otp_expiration"`
}

// RateLimitConfig holds the rate limiting configuration
type RateLimitConfig struct {
	Enabled bool          `json:"enabled"`
	Limit   int           `json:"limit"`
	Window  time.Duration `json:"window"`
}

// LoadConfig loads the configuration from a file and environment variables
func LoadConfig(filePath string) (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			Port:            "8080",
			Mode:            "debug",
			Timeout:         time.Second * 30,
			ShutdownTimeout: time.Second * 30,
			AllowedOrigins:  []string{"*"},
			RateLimit: RateLimitConfig{
				Enabled: true,
				Limit:   100,
				Window:  time.Minute,
			},
		},
		Database: DatabaseConfig{
			MaxOpenConns: 25,
			MaxIdleConns: 25,
			MaxLifetime:  time.Hour,
		},
		JWT: JWTConfig{
			AccessTokenExp:   time.Hour * 24,
			RefreshTokenExp:  time.Hour * 24 * 7,
			SigningAlgorithm: "HS256",
		},
		Validation: ValidationConfig{
			PasswordMinLength: 8,
			OTPLength:         6,
			OTPExpiration:     time.Minute * 5,
		},
	}

	// Load from file if exists
	if filePath != "" {
		file, err := os.Open(filePath)
		if err == nil {
			defer file.Close()
			if err := json.NewDecoder(file).Decode(config); err != nil {
				return nil, err
			}
		}
	}

	// Override with environment variables
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		config.Server.Mode = mode
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			config.Database.Port = port
		}
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		config.Database.Password = dbPass
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.Database = dbName
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		config.JWT.Secret = jwtSecret
	}
	if redisHost := os.Getenv("REDIS_HOST"); redisHost != "" {
		config.Redis.Host = redisHost
	}
	if redisPort := os.Getenv("REDIS_PORT"); redisPort != "" {
		if port, err := strconv.Atoi(redisPort); err == nil {
			config.Redis.Port = port
		}
	}
	if redisPass := os.Getenv("REDIS_PASSWORD"); redisPass != "" {
		config.Redis.Password = redisPass
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here
	return nil
}
