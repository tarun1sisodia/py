package controllers

// ... (Existing code) ...

// RegisterTeacher handles teacher registration.
func RegisterTeacher(c *gin.Context) {
	// ... (Existing code) ...

	// Send OTP via Firebase Authentication
	otp, err := sendFirebaseOTP(newUser.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send Firebase OTP"})
		return
	}

	// ... (Existing code) ...
}

// sendFirebaseOTP sends an OTP using Firebase Authentication.
func sendFirebaseOTP(phoneNumber string) (string, error) {
	authClient, err := GetFirebaseAuthClient(c)
	if err != nil {
		return "", err
	}

	resp, err := authClient.SendVerificationCode(context.Background(), phoneNumber, "+1") // Adjust country code as needed
	if err != nil {
		return "", err
	}

	return resp.SessionInfo.VerificationId, nil
}

// ... (Rest of the code) ...

// ResetPassword handles password reset.
func ResetPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input
	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Send password reset link using Firebase Authentication
	err := sendFirebasePasswordReset(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send password reset link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to email"})
}

// sendFirebasePasswordReset sends a password reset link using Firebase Authentication.
func sendFirebasePasswordReset(email string) error {
	authClient, err := GetFirebaseAuthClient(c)
	if err != nil {
		return err
	}

	return authClient.SendPasswordResetEmail(context.Background(), email)
}

// ... (Rest of the code) ...
