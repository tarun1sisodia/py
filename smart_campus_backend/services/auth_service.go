package services

import (
	"context"
	"fmt"
	"time"

	"smart_campus_backend/models"
	"smart_campus_backend/repositories"

	firebase "firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     repositories.UserRepository
	firebaseAuth *firebase.Client
	jwtSecret    string
}

func NewAuthService(
	userRepo repositories.UserRepository,
	firebaseAuth *firebase.Client,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		firebaseAuth: firebaseAuth,
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate role-specific fields
	if err := s.validateRegistrationFields(req); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	// Create user in Firebase
	params := (&firebase.UserToCreate{}).
		Email(req.Email).
		PhoneNumber(req.Phone).
		Password(req.Password)

	firebaseUser, err := s.firebaseAuth.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating Firebase user: %v", err)
	}

	// Create user in database
	user := &models.User{
		ID:            uuid.New().String(),
		Role:          req.Role,
		FullName:      req.FullName,
		Username:      req.Username,
		RollNumber:    req.RollNumber,
		Email:         req.Email,
		Course:        req.Course,
		AcademicYear:  req.AcademicYear,
		Phone:         req.Phone,
		HighestDegree: req.HighestDegree,
		Experience:    req.Experience,
		PasswordHash:  string(hashedPassword),
		FirebaseUID:   firebaseUser.UID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		// Rollback Firebase user creation
		_ = s.firebaseAuth.DeleteUser(ctx, firebaseUser.UID)
		return nil, fmt.Errorf("error creating user in database: %v", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %v", err)
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	var user *models.User
	var err error

	// Find user based on role
	if req.Role == models.RoleTeacher {
		user, err = s.userRepo.FindByEmail(ctx, req.Email)
	} else {
		user, err = s.userRepo.FindByRollNumber(ctx, req.RollNumber)
	}

	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %v", err)
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) VerifyPhoneNumber(ctx context.Context, userID string, phoneNumber string) error {
	// Get user from database
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Update phone number in Firebase
	params := (&firebase.UserToUpdate{}).
		PhoneNumber(phoneNumber)

	if _, err := s.firebaseAuth.UpdateUser(ctx, user.FirebaseUID, params); err != nil {
		return fmt.Errorf("error updating Firebase user: %v", err)
	}

	// Update phone verification status in database
	if err := s.userRepo.UpdatePhoneVerification(ctx, userID, true); err != nil {
		return fmt.Errorf("error updating phone verification status: %v", err)
	}

	return nil
}

func (s *AuthService) validateRegistrationFields(req *models.RegisterRequest) error {
	if req.Role == models.RoleTeacher {
		if req.Username == "" || req.Email == "" || req.HighestDegree == "" {
			return fmt.Errorf("username, email, and highest degree are required for teachers")
		}
	} else {
		if req.RollNumber == "" || req.Course == "" || req.AcademicYear == "" {
			return fmt.Errorf("roll number, course, and academic year are required for students")
		}
	}
	return nil
}

func (s *AuthService) generateTokens(user *models.User) (string, string, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
