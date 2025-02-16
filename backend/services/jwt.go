package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"smart_attendance_backend/models"
)

// GenerateToken generates a JWT token for the given user.
func GenerateToken(user models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not set in environment")
	}

	expiryStr := os.Getenv("JWT_EXPIRY_HOURS")
	expiryHours, err := strconv.Atoi(expiryStr)
	if err != nil {
		expiryHours = 24 // default expiry
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    string(user.Role),
		"exp":     time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
