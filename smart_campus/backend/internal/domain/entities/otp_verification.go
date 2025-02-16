package entities

import (
	"time"

	"github.com/google/uuid"
)

// OTPVerification represents an OTP verification record in the system
type OTPVerification struct {
	BaseEntity
	UserID       string     `json:"user_id"`
	Email        string     `json:"email"`
	OTP          string     `json:"-"` // Hide OTP in JSON responses
	Purpose      OTPPurpose `json:"purpose"`
	Status       OTPStatus  `json:"status"`
	AttemptCount int        `json:"attempt_count"`
	MaxAttempts  int        `json:"max_attempts"`
	ExpiresAt    time.Time  `json:"expires_at"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
}

// Validate validates the OTP verification entity
func (o *OTPVerification) Validate() error {
	// Add validation logic here
	return nil
}

// IsExpired checks if the OTP has expired
func (o *OTPVerification) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// IsVerified checks if the OTP has been verified
func (o *OTPVerification) IsVerified() bool {
	return o.VerifiedAt != nil
}

// MarkVerified marks the OTP as verified
func (o *OTPVerification) MarkVerified() {
	now := time.Now()
	o.VerifiedAt = &now
	o.Status = OTPStatusVerified
	o.UpdatedAt = now
}

// HasExceededMaxAttempts checks if the maximum number of attempts has been exceeded
func (o *OTPVerification) HasExceededMaxAttempts() bool {
	return o.AttemptCount >= o.MaxAttempts
}

// VerifyOTP attempts to verify the OTP
func (o *OTPVerification) VerifyOTP(providedOTP string) bool {
	if o.IsExpired() {
		o.Status = OTPStatusExpired
		o.UpdatedAt = time.Now()
		return false
	}

	if o.HasExceededMaxAttempts() {
		o.Status = OTPStatusInvalid
		o.UpdatedAt = time.Now()
		return false
	}

	o.AttemptCount++
	o.UpdatedAt = time.Now()

	if o.OTP == providedOTP {
		o.MarkVerified()
		return true
	}

	if o.HasExceededMaxAttempts() {
		o.Status = OTPStatusInvalid
	}

	return false
}

// NewOTPVerification creates a new OTP verification instance
func NewOTPVerification(userID, email, otp string, purpose OTPPurpose, expiration time.Duration) *OTPVerification {
	now := time.Now()
	return &OTPVerification{
		BaseEntity: BaseEntity{
			ID:        uuid.New().String(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:       userID,
		Email:        email,
		OTP:          otp,
		Purpose:      purpose,
		Status:       OTPStatusPending,
		AttemptCount: 0,
		MaxAttempts:  3,
		ExpiresAt:    now.Add(expiration),
	}
}
