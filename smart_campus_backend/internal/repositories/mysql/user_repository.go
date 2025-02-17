package mysql

import (
	"context"
	"database/sql"
	"errors"

	"smart_campus_backend/internal/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, role, full_name, email, phone, password_hash, firebase_uid,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Role,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.FirebaseUID,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT id, role, full_name, email, phone, password_hash, firebase_uid,
			created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, role, full_name, email, phone, password_hash, firebase_uid,
			created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &models.User{}
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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, role, full_name, email, phone, password_hash, firebase_uid,
			created_at, updated_at
		FROM users
		WHERE phone = ?
	`

	user := &models.User{}
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

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByRole(ctx context.Context, role string) ([]*models.User, error) {
	query := `
		SELECT id, role, full_name, email, phone, password_hash, firebase_uid,
			created_at, updated_at
		FROM users
		WHERE role = ?
	`

	rows, err := r.db.QueryContext(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
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
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET role = ?, full_name = ?, email = ?, phone = ?, password_hash = ?,
			firebase_uid = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		user.Role,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.FirebaseUID,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
