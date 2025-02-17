package db

import (
	"context"
	"database/sql"
	"smart_campus_backend/internal/models"
)

type UserRepository struct {
	db *Database
}

type SessionRepository struct {
	db *Database
}

type CourseRepository struct {
	db *Database
}

type DeviceRepository struct {
	db *Database
}

type AttendanceRepository struct {
	db *Database
}

// UserRepository methods
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, role, full_name, email, phone, password_hash, firebase_uid, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Role,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.FirebaseUID,
	)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `SELECT id, role, full_name, email, phone, password_hash, firebase_uid, created_at, updated_at 
			FROM users WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, role, full_name, email, phone, password_hash, firebase_uid, created_at, updated_at 
			FROM users WHERE email = ?`

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	var user models.User
	query := `SELECT id, role, full_name, email, phone, password_hash, firebase_uid, created_at, updated_at 
			FROM users WHERE phone = ?`

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&user.ID,
		&user.Role,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.FirebaseUID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users 
			SET role = ?, full_name = ?, email = ?, phone = ?, updated_at = NOW()
			WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		user.Role,
		user.FullName,
		user.Email,
		user.Phone,
		user.ID,
	)
	return err
}

// SessionRepository methods
func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	query := `INSERT INTO sessions (id, teacher_id, course_id, start_time, end_time, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.TeacherID,
		session.CourseID,
		session.StartTime,
		session.EndTime,
		session.Status,
	)
	return err
}

// CourseRepository methods
/*func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	query := `INSERT INTO courses (id, code, name, department, year, semester, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		course.ID,
		course.Code,
		course.Name,
		course.Department,
		course.Year,
		course.Semester,
	)
	return err
}

// DeviceRepository methods
func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	query := `INSERT INTO devices (id, user_id, device_id, device_name, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		device.ID,
		device.UserID,
		device.DeviceID,
		device.DeviceName,
		device.IsActive,
	)
	return err
}

// AttendanceRepository methods
func (r *AttendanceRepository) Create(ctx context.Context, attendance *models.Attendance) error {
	query := `INSERT INTO attendance (id, session_id, student_id, status, created_at, updated_at)
			VALUES (?, ?, ?, ?, NOW(), NOW())`

	_, err := r.db.ExecContext(ctx, query,
		attendance.ID,
		attendance.SessionID,
		attendance.StudentID,
		attendance.Status,
	)
	return err
}*/
