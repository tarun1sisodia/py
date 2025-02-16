package requests

// BindDeviceRequest represents the request to bind a device
type BindDeviceRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceName  string `json:"device_name" binding:"required"`
	DeviceModel string `json:"device_model" binding:"required"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	DeviceID string `json:"device_id" binding:"required"`
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Email            string `json:"email" binding:"required,email"`
	Password         string `json:"password" binding:"required,min=6"`
	FullName         string `json:"full_name" binding:"required"`
	Role             string `json:"role" binding:"required,oneof=student teacher"`
	Department       string `json:"department" binding:"required"`
	YearOfStudy      *int   `json:"year_of_study,omitempty"`
	EnrollmentNumber string `json:"enrollment_number,omitempty"`
	EmployeeID       string `json:"employee_id,omitempty"`
}

// VerifyOTPRequest represents the OTP verification request
type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// ResetPasswordRequest represents the password reset request
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
