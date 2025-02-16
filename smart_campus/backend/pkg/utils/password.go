package utils

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort    = errors.New("password too short")
	ErrPasswordTooLong     = errors.New("password too long")
	ErrPasswordNoUpper     = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower     = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber    = errors.New("password must contain at least one number")
	ErrPasswordNoSpecial   = errors.New("password must contain at least one special character")
	ErrPasswordCommonWord  = errors.New("password contains common word")
	ErrPasswordSequential  = errors.New("password contains sequential characters")
	ErrPasswordRepeating   = errors.New("password contains repeating characters")
	ErrPasswordWhitespace  = errors.New("password contains whitespace")
	ErrPasswordInvalidHash = errors.New("invalid password hash")
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72 // bcrypt's maximum input length
	BcryptCost        = 12 // Higher cost means more secure but slower
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword checks if the provided password matches the hash
func VerifyPassword(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ValidatePassword checks if a password meets security requirements
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > MaxPasswordLength {
		return ErrPasswordTooLong
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	// Check for whitespace
	if strings.ContainsAny(password, " \t\n\r") {
		return ErrPasswordWhitespace
	}

	// Check character types
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	// Check for sequential characters
	if hasSequentialChars(password) {
		return ErrPasswordSequential
	}

	// Check for repeating characters
	if hasRepeatingChars(password) {
		return ErrPasswordRepeating
	}

	return nil
}

// GenerateOTP generates a random OTP of specified length
func GenerateOTP(length int) string {
	// Calculate how many random bytes we need
	byteLength := (length * 5 / 8) + 1
	randomBytes := make([]byte, byteLength)

	// Generate random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		// In case of error, return a fallback OTP (should never happen)
		return strings.Repeat("0", length)
	}

	// Encode to base32 and take first 'length' characters
	encoded := base32.StdEncoding.EncodeToString(randomBytes)
	// Replace O and I with 8 and 9 to avoid confusion
	encoded = strings.Map(func(r rune) rune {
		switch r {
		case 'O':
			return '8'
		case 'I':
			return '9'
		default:
			return r
		}
	}, encoded)

	// Take only numeric characters and limit to desired length
	numeric := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return '0'
	}, encoded)

	if len(numeric) < length {
		// Pad with zeros if necessary (should never happen)
		numeric = numeric + strings.Repeat("0", length-len(numeric))
	}

	return numeric[:length]
}

// hasSequentialChars checks if the password contains sequential characters
func hasSequentialChars(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		if password[i]+1 == password[i+1] && password[i+1]+1 == password[i+2] {
			return true
		}
	}
	return false
}

// hasRepeatingChars checks if the password contains repeating characters
func hasRepeatingChars(password string) bool {
	if len(password) < 3 {
		return false
	}

	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

// GenerateRandomPassword generates a random password that meets all requirements
func GenerateRandomPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
	length := MinPasswordLength + 4 // A bit longer than minimum for better security

	for attempts := 0; attempts < 100; attempts++ {
		password := make([]byte, length)
		if _, err := rand.Read(password); err != nil {
			return "", err
		}

		// Map random bytes to charset
		for i := range password {
			password[i] = charset[int(password[i])%len(charset)]
		}

		// Ensure at least one of each required character type
		password[0] = 'A' // uppercase
		password[1] = 'a' // lowercase
		password[2] = '1' // number
		password[3] = '!' // special

		// Shuffle the password
		for i := len(password) - 1; i > 0; i-- {
			j := int(password[i]) % (i + 1)
			password[i], password[j] = password[j], password[i]
		}

		if err := ValidatePassword(string(password)); err == nil {
			return string(password), nil
		}
	}

	return "", errors.New("failed to generate valid password")
}
