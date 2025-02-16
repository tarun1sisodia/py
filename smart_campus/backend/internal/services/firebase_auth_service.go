package services

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type FirebaseAuthService struct {
	auth *auth.Client
}

func NewFirebaseAuthService(credentialsPath string) (*FirebaseAuthService, error) {
	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	return &FirebaseAuthService{
		auth: auth,
	}, nil
}

// VerifyPhoneNumber initiates phone number verification
func (s *FirebaseAuthService) VerifyPhoneNumber(phoneNumber string) (string, error) {
	params := (&auth.PhoneAuthParams{}).
		WithPhoneNumber(phoneNumber).
		WithRecaptchaParam("") // Handled by frontend

	verificationID, err := s.auth.VerifyPhoneNumber(context.Background(), params)
	if err != nil {
		return "", fmt.Errorf("error verifying phone number: %v", err)
	}

	return verificationID, nil
}

// VerifyIDToken verifies a Firebase ID token
func (s *FirebaseAuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := s.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token, nil
}

// GetUser gets a user by UID
func (s *FirebaseAuthService) GetUser(ctx context.Context, uid string) (*auth.UserRecord, error) {
	user, err := s.auth.GetUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

// UpdateUser updates a user's profile
func (s *FirebaseAuthService) UpdateUser(ctx context.Context, uid string, updates *auth.UserToUpdate) error {
	_, err := s.auth.UpdateUser(ctx, uid, updates)
	if err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}
	return nil
}

// RevokeTokens revokes all refresh tokens for a user
func (s *FirebaseAuthService) RevokeTokens(ctx context.Context, uid string) error {
	err := s.auth.RevokeRefreshTokens(ctx, uid)
	if err != nil {
		return fmt.Errorf("error revoking tokens: %v", err)
	}
	return nil
}

// DeleteUser deletes a user
func (s *FirebaseAuthService) DeleteUser(ctx context.Context, uid string) error {
	err := s.auth.DeleteUser(ctx, uid)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}
