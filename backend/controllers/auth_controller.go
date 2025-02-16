// Teacher registration, login, and OTP verification controller for Smart Attendance App

package controllers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"smart_attendance_backend/models"
	"smart_attendance_backend/services"
)

// TeacherRegistrationRequest holds the request payload for teacher registration
// (Reapplying struct with an extra newline after it to resolve potential formatting issues)
type TeacherRegistrationRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Phone           string `json:"phone" binding:"required"`
	HighestDegree   string `json:"highest_degree" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
	Experience      string `json:"experience" binding:"required"`
}

// TeacherRegister handles teacher registration
func TeacherRegister(c *gin.Context) {
	var req TeacherRegistrationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if password and confirm match
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check for duplicate username or phone
	var existing models.User
	if err := models.GetDB().Where("username = ? OR phone = ?", req.Username, req.Phone).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or phone number already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new user record
	newUser := models.User{
		Role:          models.RoleTeacher,
		FullName:      req.FullName,
		Username:      &req.Username,
		Email:         &req.Email,
		Phone:         req.Phone,
		HighestDegree: &req.HighestDegree,
		Experience:    &req.Experience,
		PasswordHash:  string(hashedPassword),
	}

	if err := models.GetDB().Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate a 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(100000 + rand.Intn(900000))

	// Create OTP verification record
	otpRec := models.OTPVerification{
		UserID:    newUser.ID,
		OTPCode:   otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := models.GetDB().Create(&otpRec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP record"})
		return
	}

	// Send OTP via Twilio (stub implementation)
	if err := services.SendSMS(req.Phone, "Your OTP code is: "+otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to phone number. Use this user ID for OTP verification.", "user_id": newUser.ID})
}

// New structs for additional authentication endpoints

type StudentRegistrationRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	RollNumber      string `json:"roll_number" binding:"required"`
	Course          string `json:"course" binding:"required"`
	AcademicYear    string `json:"academic_year" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=6"`
}

type TeacherLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type StudentLoginRequest struct {
	RollNumber string `json:"roll_number" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type ResetPasswordRequest struct {
	UserID             string `json:"user_id" binding:"required"`
	OTP                string `json:"otp" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required,min=6"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required,min=6"`
}

// StudentRegister handles student registration
func StudentRegister(c *gin.Context) {
	var req StudentRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Check for duplicate roll number or phone
	var existing models.User
	if err := models.GetDB().Where("roll_number = ? OR phone = ?", req.RollNumber, req.Phone).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Roll number or phone number already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create new student user record
	newUser := models.User{
		Role:         models.RoleStudent,
		FullName:     req.FullName,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
	}
	// Set student-specific fields
	newUser.RollNumber = &req.RollNumber
	newUser.Course = &req.Course
	newUser.AcademicYear = &req.AcademicYear

	if err := models.GetDB().Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate a 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(100000 + rand.Intn(900000))

	// Create OTP verification record
	otpRec := models.OTPVerification{
		UserID:    newUser.ID,
		OTPCode:   otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := models.GetDB().Create(&otpRec).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create OTP record"})
		return
	}

	// Send OTP via Twilio (stub implementation)
	if err := services.SendSMS(req.Phone, "Your OTP code is: "+otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to phone number. Use this user ID for OTP verification.", "user_id": newUser.ID})
}

// TeacherLogin handles teacher login
func TeacherLogin(c *gin.Context) {
	var req TeacherLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.GetDB().Where("email = ? AND role = ?", req.Email, models.RoleTeacher).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := services.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// StudentLogin handles student login
func StudentLogin(c *gin.Context) {
	var req StudentLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := models.GetDB().Where("roll_number = ? AND role = ?", req.RollNumber, models.RoleStudent).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := services.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

// ResetPassword handles password reset via OTP
func ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.NewPassword != req.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	db := models.GetDB()
	var otpRecord models.OTPVerification
	if err := db.Where("user_id = ? AND verified = ?", req.UserID, false).Order("created_at desc").First(&otpRecord).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP record not found"})
		return
	}

	if otpRecord.IsExpired() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired. Please request a new one."})
		return
	}

	if otpRecord.OTPCode != req.OTP {
		otpRecord.IncrementAttempts()
		db.Save(&otpRecord)
		if otpRecord.HasExceededMaxAttempts(3) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP verification failed, maximum attempts exceeded."})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP."})
		}
		return
	}

	// Update user's password
	var user models.User
	if err := db.Where("id = ?", req.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}
	user.PasswordHash = string(hashedPassword)
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	otpRecord.MarkAsVerified()
	db.Save(&otpRecord)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
