package services

import (
	"context"
	"smart_campus_backend/internal/models"

	firebase "firebase.google.com/go/v4/auth"
)

// Services holds all service instances
type Services struct {
	AuthService *AuthService
}

// Repositories holds all repository instances
type Repositories struct {
	UserRepo UserRepository
}

// Repository interfaces
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByPhone(ctx context.Context, phone string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

func NewAuthService(userRepo UserRepository, firebaseAuth *firebase.Client) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		firebaseAuth: firebaseAuth,
	}
}
