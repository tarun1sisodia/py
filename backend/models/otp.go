package models

import "time"

// OTPVerification represents an OTP verification record.
type OTPVerification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	OTPCode   string    `json:"otp_code"`
	ExpiresAt time.Time `json:"expires_at"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
}
