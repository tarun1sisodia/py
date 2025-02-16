package repositories

import (
	"context"
	"smart_campus_backend/models"
)

type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *models.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*models.User, error)

	// FindByEmail finds a teacher by email
	FindByEmail(ctx context.Context, email string) (*models.User, error)

	// FindByRollNumber finds a student by roll number
	FindByRollNumber(ctx context.Context, rollNumber string) (*models.User, error)

	// FindByFirebaseUID finds a user by Firebase UID
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*models.User, error)

	// Update updates a user
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id string) error

	// UpdatePhoneVerification updates the phone verification status
	UpdatePhoneVerification(ctx context.Context, userID string, verified bool) error
}
