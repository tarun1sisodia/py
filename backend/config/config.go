package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config stores the application configuration.
type Config struct {
	Port            int
	DatabaseURL     string
	FirebaseProjectID string
	JWTSecretKey    string
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() *Config {
	// Load .env file if it exists.
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: ", err)
	}

	// Get the port from the environment variables, default to 8080
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
		log.Println("Defaulting to port 8080")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	// Get the database URL from the environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Get the Firebase Project ID from the environment variables
	firebaseProjectID := os.Getenv("FIREBASE_PROJECT_ID")
	if firebaseProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID environment variable not set")
	}

	// Get the JWT secret key from the environment variables
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Println("JWT_SECRET_KEY environment variable not set, using default secret key")
		jwtSecretKey = "secret" // Use a default secret key if not set in environment variables
	}

	return &Config{
		Port:            port,
		DatabaseURL:     databaseURL,
		FirebaseProjectID: firebaseProjectID,
		JWTSecretKey:    jwtSecretKey,
	}
}
