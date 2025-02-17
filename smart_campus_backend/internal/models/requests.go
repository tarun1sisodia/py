package models

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Role     string `json:"role" binding:"required,oneof=student teacher"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents the token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// VerifyPhoneRequest represents the phone verification request
type VerifyPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// UpdateProfileRequest represents the profile update request
type UpdateProfileRequest struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Phone    string `json:"phone,omitempty"`
}
