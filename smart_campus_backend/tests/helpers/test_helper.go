package helpers

import (
	"database/sql"
	"fmt"
	"log"

	"smart_campus_backend/internal/models"
	"smart_campus_backend/tests/config"
)

var testDB *sql.DB

// InitTestDB initializes the test database
func InitTestDB() *sql.DB {
	cfg := config.LoadTestConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping test database: %v", err)
	}

	testDB = db
	return db
}

// CleanupTestDB cleans up test data after tests
func CleanupTestDB() {
	if testDB != nil {
		tables := []string{"users", "devices"}
		for _, table := range tables {
			_, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table))
			if err != nil {
				log.Printf("Failed to truncate table %s: %v", table, err)
			}
		}
	}
}

// CreateTestUser creates a test user in the database
func CreateTestUser() (*models.User, error) {
	user := &models.User{
		Role:         "student",
		FullName:     "Test User",
		Email:        "test@example.com",
		Phone:        "+1234567890",
		PasswordHash: "$2a$10$abcdefghijklmnopqrstuvwxyz123456",
		FirebaseUID:  "test-firebase-uid",
	}

	query := `
		INSERT INTO users (role, full_name, email, phone, password_hash, firebase_uid)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := testDB.Exec(query,
		user.Role,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.FirebaseUID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create test user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %v", err)
	}

	user.ID = fmt.Sprintf("%d", id)
	return user, nil
}

// GetTestUser retrieves a test user from the database
func GetTestUser(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, role, full_name, email, phone, firebase_uid, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	err := testDB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get test user: %v", err)
	}

	return user, nil
}
