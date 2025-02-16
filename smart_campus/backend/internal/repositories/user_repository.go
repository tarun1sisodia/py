package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/database"
	"smart_campus/internal/models"

	"github.com/google/uuid"
)

type MySQLUserRepository struct {
	db *database.MySQLDB
}

func NewMySQLUserRepository(db *database.MySQLDB) models.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *models.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (
			id, role, email, password_hash, full_name, enrollment_number,
			employee_id, department, year_of_study, device_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		user.ID, user.Role, user.Email, user.PasswordHash, user.FullName,
		user.EnrollmentNumber, user.EmployeeID, user.Department,
		user.YearOfStudy, user.DeviceID, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (r *MySQLUserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, role, email, password_hash, full_name, enrollment_number,
		employee_id, department, year_of_study, device_id, created_at, updated_at
		FROM users WHERE id = ?
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Role, &user.Email, &user.PasswordHash, &user.FullName,
		&user.EnrollmentNumber, &user.EmployeeID, &user.Department,
		&user.YearOfStudy, &user.DeviceID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %v", err)
	}

	return user, nil
}

func (r *MySQLUserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, role, email, password_hash, full_name, enrollment_number,
		employee_id, department, year_of_study, device_id, created_at, updated_at
		FROM users WHERE email = ?
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Role, &user.Email, &user.PasswordHash, &user.FullName,
		&user.EnrollmentNumber, &user.EmployeeID, &user.Department,
		&user.YearOfStudy, &user.DeviceID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by email: %v", err)
	}

	return user, nil
}

func (r *MySQLUserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users SET
			role = ?, email = ?, full_name = ?, enrollment_number = ?,
			employee_id = ?, department = ?, year_of_study = ?, device_id = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		user.Role, user.Email, user.FullName, user.EnrollmentNumber,
		user.EmployeeID, user.Department, user.YearOfStudy, user.DeviceID,
		user.UpdatedAt, user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MySQLUserRepository) Delete(id string) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MySQLUserRepository) List(offset, limit int) ([]*models.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name, enrollment_number,
		employee_id, department, year_of_study, device_id, created_at, updated_at
		FROM users LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Role, &user.Email, &user.PasswordHash, &user.FullName,
			&user.EnrollmentNumber, &user.EmployeeID, &user.Department,
			&user.YearOfStudy, &user.DeviceID, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

func (r *MySQLUserRepository) UpdatePassword(id string, passwordHash string) error {
	query := `
		UPDATE users SET
			password_hash = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, passwordHash, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MySQLUserRepository) UpdateDeviceID(id string, deviceID *string) error {
	query := `
		UPDATE users SET
			device_id = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, deviceID, time.Now(), id)
	if err != nil {
		return fmt.Errorf("error updating device ID: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
