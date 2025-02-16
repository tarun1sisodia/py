package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/domain"
	"smart_campus/internal/domain/entities"
	"smart_campus/internal/domain/repositories"
)

// OTPVerificationRepository implements the repositories.OTPVerificationRepository interface
type OTPVerificationRepository struct {
	conn *Connection
}

// NewOTPVerificationRepository creates a new MySQL OTP verification repository
func NewOTPVerificationRepository(conn *Connection) repositories.OTPVerificationRepository {
	return &OTPVerificationRepository{conn: conn}
}

// Create creates a new OTP verification record
func (r *OTPVerificationRepository) Create(ctx context.Context, verification *entities.OTPVerification) error {
	query := `
		INSERT INTO otp_verifications (
			id, user_id, otp, purpose, status, attempt_count,
			max_attempts, expires_at, verified_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.DB().ExecContext(ctx, query,
		verification.ID, verification.UserID, verification.OTP,
		verification.Purpose, verification.Status, verification.AttemptCount,
		verification.MaxAttempts, verification.ExpiresAt, verification.VerifiedAt,
		verification.CreatedAt, verification.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating OTP verification: %v", err)
	}

	return nil
}

// GetByID retrieves an OTP verification record by ID
func (r *OTPVerificationRepository) GetByID(ctx context.Context, id string) (*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications WHERE id = ?
	`

	verification := &entities.OTPVerification{}
	err := r.conn.DB().QueryRowContext(ctx, query, id).Scan(
		&verification.ID, &verification.UserID, &verification.OTP,
		&verification.Purpose, &verification.Status, &verification.AttemptCount,
		&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
		&verification.CreatedAt, &verification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting OTP verification: %v", err)
	}

	return verification, nil
}

// Update updates an existing OTP verification record
func (r *OTPVerificationRepository) Update(ctx context.Context, verification *entities.OTPVerification) error {
	query := `
		UPDATE otp_verifications SET
			status = ?, attempt_count = ?, verified_at = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		verification.Status, verification.AttemptCount,
		verification.VerifiedAt, verification.UpdatedAt,
		verification.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating OTP verification: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Delete deletes an OTP verification record by ID
func (r *OTPVerificationRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM otp_verifications WHERE id = ?"

	result, err := r.conn.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting OTP verification: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetByUser retrieves OTP verification records for a user
func (r *OTPVerificationRepository) GetByUser(ctx context.Context, userID string) ([]*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user OTP verifications: %v", err)
	}
	defer rows.Close()

	var verifications []*entities.OTPVerification
	for rows.Next() {
		verification := &entities.OTPVerification{}
		err := rows.Scan(
			&verification.ID, &verification.UserID, &verification.OTP,
			&verification.Purpose, &verification.Status, &verification.AttemptCount,
			&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
			&verification.CreatedAt, &verification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OTP verification row: %v", err)
		}
		verifications = append(verifications, verification)
	}

	return verifications, nil
}

// GetActiveByUser retrieves active OTP verification records for a user
func (r *OTPVerificationRepository) GetActiveByUser(ctx context.Context, userID string) ([]*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications
		WHERE user_id = ? AND status = ? AND expires_at > ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID, entities.OTPStatusPending, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error getting active OTP verifications: %v", err)
	}
	defer rows.Close()

	var verifications []*entities.OTPVerification
	for rows.Next() {
		verification := &entities.OTPVerification{}
		err := rows.Scan(
			&verification.ID, &verification.UserID, &verification.OTP,
			&verification.Purpose, &verification.Status, &verification.AttemptCount,
			&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
			&verification.CreatedAt, &verification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OTP verification row: %v", err)
		}
		verifications = append(verifications, verification)
	}

	return verifications, nil
}

// GetByUserAndPurpose retrieves OTP verification records for a user with a specific purpose
func (r *OTPVerificationRepository) GetByUserAndPurpose(ctx context.Context, userID string, purpose entities.OTPPurpose) ([]*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications
		WHERE user_id = ? AND purpose = ?
		ORDER BY created_at DESC
	`

	rows, err := r.conn.DB().QueryContext(ctx, query, userID, purpose)
	if err != nil {
		return nil, fmt.Errorf("error getting OTP verifications by purpose: %v", err)
	}
	defer rows.Close()

	var verifications []*entities.OTPVerification
	for rows.Next() {
		verification := &entities.OTPVerification{}
		err := rows.Scan(
			&verification.ID, &verification.UserID, &verification.OTP,
			&verification.Purpose, &verification.Status, &verification.AttemptCount,
			&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
			&verification.CreatedAt, &verification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning OTP verification row: %v", err)
		}
		verifications = append(verifications, verification)
	}

	return verifications, nil
}

// GetLatestByUser retrieves the latest OTP verification record for a user
func (r *OTPVerificationRepository) GetLatestByUser(ctx context.Context, userID string) (*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	verification := &entities.OTPVerification{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID).Scan(
		&verification.ID, &verification.UserID, &verification.OTP,
		&verification.Purpose, &verification.Status, &verification.AttemptCount,
		&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
		&verification.CreatedAt, &verification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting latest OTP verification: %v", err)
	}

	return verification, nil
}

// VerifyOTP verifies an OTP code for a user
func (r *OTPVerificationRepository) VerifyOTP(ctx context.Context, userID string, otp string) (*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications
		WHERE user_id = ? AND otp = ? AND status = ? AND expires_at > ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	verification := &entities.OTPVerification{}
	err := r.conn.DB().QueryRowContext(ctx, query, userID, otp, entities.OTPStatusPending, time.Now()).Scan(
		&verification.ID, &verification.UserID, &verification.OTP,
		&verification.Purpose, &verification.Status, &verification.AttemptCount,
		&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
		&verification.CreatedAt, &verification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("error verifying OTP: %v", err)
	}

	return verification, nil
}

// InvalidateOTP marks an OTP verification record as invalid
func (r *OTPVerificationRepository) InvalidateOTP(ctx context.Context, id string) error {
	query := `
		UPDATE otp_verifications
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.conn.DB().ExecContext(ctx, query,
		entities.OTPStatusInvalid, time.Now(), id,
	)

	if err != nil {
		return fmt.Errorf("error invalidating OTP: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// DeleteExpired deletes expired OTP verification records
func (r *OTPVerificationRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	query := "DELETE FROM otp_verifications WHERE expires_at < ?"

	_, err := r.conn.DB().ExecContext(ctx, query, before)
	if err != nil {
		return fmt.Errorf("error deleting expired OTP verifications: %v", err)
	}

	return nil
}

// CountActiveOTPs counts the number of active OTP verification records for a user
func (r *OTPVerificationRepository) CountActiveOTPs(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM otp_verifications
		WHERE user_id = ? AND status = ? AND expires_at > ?
	`

	var count int
	err := r.conn.DB().QueryRowContext(ctx, query, userID, entities.OTPStatusPending, time.Now()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting active OTPs: %v", err)
	}

	return count, nil
}

// GetByOTP retrieves an OTP verification record by the OTP code
func (r *OTPVerificationRepository) GetByOTP(ctx context.Context, otp string) (*entities.OTPVerification, error) {
	query := `
		SELECT id, user_id, otp, purpose, status, attempt_count,
		       max_attempts, expires_at, verified_at, created_at, updated_at
		FROM otp_verifications
		WHERE otp = ? AND status = ? AND expires_at > ?
		LIMIT 1
	`

	verification := &entities.OTPVerification{}
	err := r.conn.DB().QueryRowContext(ctx, query, otp, entities.OTPStatusPending, time.Now()).Scan(
		&verification.ID, &verification.UserID, &verification.OTP,
		&verification.Purpose, &verification.Status, &verification.AttemptCount,
		&verification.MaxAttempts, &verification.ExpiresAt, &verification.VerifiedAt,
		&verification.CreatedAt, &verification.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error getting OTP verification by code: %v", err)
	}

	return verification, nil
}
