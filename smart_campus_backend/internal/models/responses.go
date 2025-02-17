package models

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PhoneVerificationResponse represents the response for phone verification
type PhoneVerificationResponse struct {
	VerificationID string `json:"verification_id"`
	ExpiresAt      int64  `json:"expires_at"`
}
