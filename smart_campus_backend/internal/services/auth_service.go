package services

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"smart_campus_backend/internal/models"
)

type AuthService struct {
	userRepo     UserRepository
	firebaseAuth *auth.Client
}

// InitAuthService initializes a new AuthService
func InitAuthService(userRepo UserRepository, app *firebase.App) (*AuthService, error) {
	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	return &AuthService{
		userRepo:     userRepo,
		firebaseAuth: auth,
	}, nil
}

func (s *AuthService) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
		return
	}

	// Create user in Firebase
	params := (&auth.UserToCreate{}).
		Email(req.Email).
		PhoneNumber(req.Phone).
		Password(req.Password)

	firebaseUser, err := s.firebaseAuth.CreateUser(c.Request.Context(), params)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error creating Firebase user: %v", err)})
		return
	}

	// Create user in database
	user := &models.User{
		ID:           uuid.New().String(),
		Role:         req.Role,
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashedPassword),
		FirebaseUID:  firebaseUser.UID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(c.Request.Context(), user); err != nil {
		// Rollback Firebase user creation
		_ = s.firebaseAuth.DeleteUser(c.Request.Context(), firebaseUser.UID)
		c.JSON(500, gin.H{"error": fmt.Sprintf("Error creating user in database: %v", err)})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating tokens"})
		return
	}

	c.JSON(201, models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthService) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := s.userRepo.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating tokens"})
		return
	}

	c.JSON(200, models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *AuthService) VerifyPhone(c *gin.Context) {
	var req models.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Generate verification code
	verificationID := uuid.New().String()
	// In a real implementation, you would send this code via SMS
	// For now, we'll just return it

	c.JSON(200, models.PhoneVerificationResponse{
		VerificationID: verificationID,
		ExpiresAt:      time.Now().Add(10 * time.Minute).Unix(),
	})
}

func (s *AuthService) RefreshToken(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	user, err := s.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating tokens"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *AuthService) Logout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	user, err := s.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	// Revoke Firebase tokens
	if err := s.firebaseAuth.RevokeRefreshTokens(c.Request.Context(), user.FirebaseUID); err != nil {
		c.JSON(500, gin.H{"error": "Error revoking tokens"})
		return
	}

	c.JSON(200, gin.H{"message": "Successfully logged out"})
}

func (s *AuthService) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	user, err := s.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func (s *AuthService) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "User not found in context"})
		return
	}

	var req struct {
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := s.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	user.FullName = req.FullName
	user.Phone = req.Phone
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(c.Request.Context(), user); err != nil {
		c.JSON(500, gin.H{"error": "Error updating user"})
		return
	}

	c.JSON(200, user)
}

func (s *AuthService) generateTokens(user *models.User) (string, string, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte("your-refresh-secret-key"))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
