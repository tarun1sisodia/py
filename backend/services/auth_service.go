package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"firebase.google.com/go/auth"
	"github.com/golang-jwt/jwt"
	"smart-attendance/config"
	"smart-attendance/models"
)

// AuthService handles authentication logic.
type AuthService struct {
	// Add any dependencies here, such as a database repository.
}

// NewAuthService creates a new AuthService.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// CreateUser creates a new user.
func CreateUser(user *models.User) error {
	_, err := config.DB.Exec("INSERT INTO users (id, role, full_name, username, roll_number, email, course, academic_year, phone, highest_degree, experience, password_hash, created_at, updated_at, firebase_uid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.Role, user.FullName, user.Username, user.RollNumber, user.Email, user.Course, user.AcademicYear, user.Phone, user.HighestDegree, user.Experience, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.FirebaseUID)
	return err
}

// GetUserByID gets a user by ID.
func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	err := config.DB.QueryRow("SELECT id, role, full_name, username, roll_number, email, course, academic_year, phone, highest_degree, experience, password_hash, created_at, updated_at, firebase_uid FROM users WHERE id = ?", id).Scan(
		&user.ID, &user.Role, &user.FullName, &user.Username, &user.RollNumber, &user.Email, &user.Course, &user.AcademicYear, &user.Phone, &user.HighestDegree, &user.Experience, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByEmail gets a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := config.DB.QueryRow("SELECT id, role, full_name, username, roll_number, email, course, academic_year, phone, highest_degree, experience, password_hash, created_at, updated_at, firebase_uid FROM users WHERE email = ?", email).Scan(
		&user.ID, &user.Role, &user.FullName, &user.Username, &user.RollNumber, &user.Email, &user.Course, &user.AcademicYear, &user.Phone, &user.HighestDegree, &user.Experience, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByRollNumber gets a user by roll number.
func GetUserByRollNumber(rollNumber string) (*models.User, error) {
	user := &models.User{}
	err := config.DB.QueryRow("SELECT id, role, full_name, username, roll_number, email, course, academic_year, phone, highest_degree, experience, password_hash, created_at, updated_at, firebase_uid FROM users WHERE roll_number = ?", rollNumber).Scan(
		&user.ID, &user.Role, &user.FullName, &user.Username, &user.RollNumber, &user.Email, &user.Course, &user.AcademicYear, &user.Phone, &user.HighestDegree, &user.Experience, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.FirebaseUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser updates a user.
func UpdateUser(user *models.User) error {
	_, err := config.DB.Exec("UPDATE users SET role = ?, full_name = ?, username = ?, roll_number = ?, email = ?, course = ?, academic_year = ?, phone = ?, highest_degree = ?, experience = ?, password_hash = ?, created_at = ?, updated_at = ?, firebase_uid = ? WHERE id = ?",
		user.Role, user.FullName, user.Username, user.RollNumber, user.Email, user.Course, user.AcademicYear, user.Phone, user.HighestDegree, user.Experience, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.FirebaseUID, user.ID)
	return err
}

// DeleteUser deletes a user.
func DeleteUser(id string) error {
	_, err := config.DB.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// GenerateJWT generates a JWT token.
func GenerateJWT(userID string) (string, error) {
	// Set the token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "secret" // Use a default secret key if not set in environment variables
		log.Println("JWT_SECRET_KEY environment variable not set, using default secret key")
	}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetFirebaseAuthClient initializes and returns a Firebase Auth client.
func GetFirebaseAuthClient(ctx context.Context, projectID string) (*auth.Client, error) {
	// Initialize Firebase Admin SDK
	opt := option.NewCredentialsFile("path/to/your/serviceAccountKey.json") // Replace with the actual path to your service account key file
	config := &firebase.Config{
		ProjectID: projectID,
	}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Printf("Failed to initialize Firebase app: %v\n", err)
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	// Initialize Firebase Auth client
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("Failed to initialize Firebase Auth client: %v\n", err)
		return nil, fmt.Errorf("error initializing auth: %v", err)
	}

	return client, nil
}
