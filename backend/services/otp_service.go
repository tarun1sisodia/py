package services

import (
	"log"
	"time"

	"github.com/google/uuid"
	"smart-attendance/config"
	"smart-attendance/models"
	"smart-attendance/utils"
)

// GenerateOTP generates a new OTP code and stores it in the database.
func GenerateOTP(userID string) (string, error) {
	// Generate a random OTP code
	otpCode := utils.GenerateOTP()

	// Set the expiration time for the OTP (e.g., 5 minutes)
	expiresAt := time.Now().Add(time.Minute * 5)

	// Create a new OTP verification record
	newOTP := models.OTPVerification{
		ID:        uuid.New().String(),
		UserID:    userID,
		OTPCode:   otpCode,
		ExpiresAt: expiresAt,
		Verified:  false,
		CreatedAt: time.Now(),
	}

	// Save the OTP verification record to the database
	err := CreateOTPVerification(&newOTP)
	if err != nil {
		log.Println("Error creating OTP verification record:", err)
		return "", err
	}

	return otpCode, nil
}

// ValidateOTP validates the OTP code for a given user.
func ValidateOTP(userID, otpCode string) (bool, error) {
	// Retrieve the OTP verification record from the database
	otpVerification, err := GetLatestOTPVerification(userID)
	if err != nil {
		log.Println("Error retrieving OTP verification record:", err)
		return false, err
	}

	// Check if the OTP verification record exists and is not expired
	if otpVerification == nil || otpVerification.ExpiresAt.Before(time.Now()) || otpVerification.Verified {
		return false, nil
	}

	// Check if the OTP code matches
	if otpVerification.OTPCode != otpCode {
		return false, nil
	}

	// Mark the OTP as verified
	otpVerification.Verified = true
	err = UpdateOTPVerification(otpVerification)
	if err != nil {
		log.Println("Error updating OTP verification record:", err)
		return false, err
	}

	return true, nil
}

// CreateOTPVerification creates a new OTP verification record.
func CreateOTPVerification(otpVerification *models.OTPVerification) error {
	_, err := config.DB.Exec("INSERT INTO otp_verifications (id, user_id, otp_code, expires_at, verified, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		otpVerification.ID, otpVerification.UserID, otpVerification.OTPCode, otpVerification.ExpiresAt, otpVerification.Verified, otpVerification.CreatedAt)
	return err
}

// GetLatestOTPVerification retrieves the latest OTP verification record for a user.
func GetLatestOTPVerification(userID string) (*models.OTPVerification, error) {
	otpVerification := &models.OTPVerification{}
	err := config.DB.QueryRow("SELECT id, user_id, otp_code, expires_at, verified, created_at FROM otp_verifications WHERE user_id = ? ORDER BY created_at DESC LIMIT 1", userID).Scan(
		&otpVerification.ID, &otpVerification.UserID, &otpVerification.OTPCode, &otpVerification.ExpiresAt, &otpVerification.Verified, &otpVerification.CreatedAt)
	if err != nil {
		return nil, err
	}
	return otpVerification, nil
}

// UpdateOTPVerification updates an OTP verification record.
func UpdateOTPVerification(otpVerification *models.OTPVerification) error {
	_, err := config.DB.Exec("UPDATE otp_verifications SET user_id = ?, otp_code = ?, expires_at = ?, verified = ?, created_at = ? WHERE id = ?",
		otpVerification.UserID, otpVerification.OTPCode, otpVerification.ExpiresAt, otpVerification.Verified, otpVerification.CreatedAt, otpVerification.ID)
	return err
}
