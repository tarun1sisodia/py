package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
)

// UserRepository implements the repositories.UserRepository interface
type UserRepository struct {
	conn *Connection
}

// NewUserRepository creates a new MySQL user repository
func NewUserRepository(conn *Connection) repositories.UserRepository {
	return &UserRepository{conn: conn}
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name,
		       enrollment_number, employee_id, department,
		       year_of_study, device_id, created_at,
		       updated_at, last_login, is_active
		FROM users
		WHERE id = ?
	`

	user := &entities.User{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Role, &user.Email, &user.PasswordHash,
		&user.Name, &user.EnrollmentNumber, &user.EmployeeID,
		&user.Department, &user.YearOfStudy, &user.DeviceID,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
		&user.Active,
	)

	if err == sql.ErrNoRows {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name,
		       enrollment_number, employee_id, department,
		       year_of_study, device_id, created_at,
		       updated_at, last_login, is_active
		FROM users
		WHERE email = ?
	`

	user := &entities.User{}
	err := r.conn.DB().QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Role, &user.Email, &user.PasswordHash,
		&user.Name, &user.EnrollmentNumber, &user.EmployeeID,
		&user.Department, &user.YearOfStudy, &user.DeviceID,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
		&user.Active,
	)

	if err == sql.ErrNoRows {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by email: %w", err)
	}

	return user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			id, role, email, password_hash, full_name,
			enrollment_number, employee_id, department,
			year_of_study, device_id, created_at,
			updated_at, last_login, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		user.ID, user.Role, user.Email, user.PasswordHash,
		user.Name, user.EnrollmentNumber, user.EmployeeID,
		user.Department, user.YearOfStudy, user.DeviceID,
		user.CreatedAt, user.UpdatedAt, user.LastLogin,
		user.Active,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return repositories.ErrDuplicate
		}
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users SET
			role = ?, email = ?, password_hash = ?,
			full_name = ?, enrollment_number = ?,
			employee_id = ?, department = ?,
			year_of_study = ?, device_id = ?,
			updated_at = ?, last_login = ?,
			is_active = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		user.Role, user.Email, user.PasswordHash,
		user.Name, user.EnrollmentNumber, user.EmployeeID,
		user.Department, user.YearOfStudy, user.DeviceID,
		user.UpdatedAt, user.LastLogin, user.Active,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// List lists all users with pagination
func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name,
		       enrollment_number, employee_id, department,
		       year_of_study, device_id, created_at,
		       updated_at, last_login, is_active
		FROM users
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Role, &user.Email, &user.PasswordHash,
			&user.Name, &user.EnrollmentNumber, &user.EmployeeID,
			&user.Department, &user.YearOfStudy, &user.DeviceID,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
			&user.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Count returns the total number of users
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.conn.DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting users: %w", err)
	}
	return count, nil
}

// GetByDeviceID retrieves a user by device ID
func (r *UserRepository) GetByDeviceID(ctx context.Context, deviceID string) (*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name,
		       enrollment_number, employee_id, department,
		       year_of_study, device_id, created_at,
		       updated_at, last_login, is_active
		FROM users
		WHERE device_id = ?
	`

	user := &entities.User{}
	err := r.conn.DB().QueryRowContext(ctx, query, deviceID).Scan(
		&user.ID, &user.Role, &user.Email, &user.PasswordHash,
		&user.Name, &user.EnrollmentNumber, &user.EmployeeID,
		&user.Department, &user.YearOfStudy, &user.DeviceID,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
		&user.Active,
	)

	if err == sql.ErrNoRows {
		return nil, repositories.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding user by device ID: %w", err)
	}

	return user, nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `
		UPDATE users SET
			password_hash = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		passwordHash,
		time.Now(),
		userID,
	)

	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// UpdateDeviceID updates a user's device ID
func (r *UserRepository) UpdateDeviceID(ctx context.Context, userID, deviceID string) error {
	query := `
		UPDATE users SET
			device_id = ?,
			updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		deviceID,
		time.Now(),
		userID,
	)

	if err != nil {
		return fmt.Errorf("error updating device ID: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// GetStudentsByDepartmentAndYear retrieves students by department and year
func (r *UserRepository) GetStudentsByDepartmentAndYear(ctx context.Context, department string, year int) ([]*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, name, enrollment_number, employee_id, 
		department, year_of_study, device_id, last_login, is_active, created_at, updated_at
		FROM users 
		WHERE role = ? AND department = ? AND year_of_study = ?
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, entities.UserRoleStudent, department, year)
	if err != nil {
		return nil, fmt.Errorf("error getting students: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		var lastLogin sql.NullTime
		err := rows.Scan(
			&user.ID,
			&user.Role,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.EnrollmentNumber,
			&user.EmployeeID,
			&user.Department,
			&user.YearOfStudy,
			&user.DeviceID,
			&lastLogin,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning student: %w", err)
		}
		if lastLogin.Valid {
			user.LastLogin = lastLogin.Time
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating students: %w", err)
	}

	return users, nil
}

// GetTeachersByDepartment retrieves teachers by department
func (r *UserRepository) GetTeachersByDepartment(ctx context.Context, department string) ([]*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, name, enrollment_number, employee_id, 
		department, year_of_study, device_id, last_login, is_active, created_at, updated_at
		FROM users 
		WHERE role = ? AND department = ?
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, entities.UserRoleTeacher, department)
	if err != nil {
		return nil, fmt.Errorf("error getting teachers: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		var lastLogin sql.NullTime
		err := rows.Scan(
			&user.ID,
			&user.Role,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.EnrollmentNumber,
			&user.EmployeeID,
			&user.Department,
			&user.YearOfStudy,
			&user.DeviceID,
			&lastLogin,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning teacher: %w", err)
		}
		if lastLogin.Valid {
			user.LastLogin = lastLogin.Time
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating teachers: %w", err)
	}

	return users, nil
}

// UpdateStatus updates a user's status
func (r *UserRepository) UpdateStatus(ctx context.Context, userID string, active bool) error {
	query := `
		UPDATE users
		SET is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query, active, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("error updating user status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// UpdateLastLogin updates a user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET last_login = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	result, err := r.conn.DB().ExecContext(ctx, query, now, now, userID)
	if err != nil {
		return fmt.Errorf("error updating last login: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rows == 0 {
		return repositories.ErrNotFound
	}

	return nil
}

// FindByRole finds users by role
func (r *UserRepository) FindByRole(ctx context.Context, role entities.UserRole) ([]*entities.User, error) {
	query := `
		SELECT id, role, email, password_hash, full_name,
		       enrollment_number, employee_id, department,
		       year_of_study, device_id, created_at,
		       updated_at, last_login, is_active
		FROM users
		WHERE role = ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("error finding users by role: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		err := rows.Scan(
			&user.ID, &user.Role, &user.Email, &user.PasswordHash,
			&user.Name, &user.EnrollmentNumber, &user.EmployeeID,
			&user.Department, &user.YearOfStudy, &user.DeviceID,
			&user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
			&user.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
