package config

import (
	"os"
	"smart_campus_backend/config"
)

func LoadTestConfig() *config.Config {
	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASSWORD", "root1234")
	os.Setenv("DB_NAME", "smart_campus_test")
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Setenv("FIREBASE_CREDENTIALS_PATH", "../config/firebase-credentials-test.json")

	// Load config with test values
	cfg, _ := config.Load()
	return cfg
}
