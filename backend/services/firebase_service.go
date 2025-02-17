package services

// ... (Existing code) ...

// VerifyFirebaseToken verifies the Firebase ID token.
func VerifyFirebaseToken(ctx context.Context, client *auth.Client, idToken string) (*auth.Token, error) {
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token, nil
}

// GetUserByFirebaseUID gets a user by Firebase UID.
func GetUserByFirebaseUID(firebaseUID string) (*auth.UserRecord, error) {
	ctx := context.Background()
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	client, err := NewFirebaseAuthClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing auth: %v", err)
	}
	user, err := client.GetUser(ctx, firebaseUID)
	if err != nil {
		return nil, fmt.Errorf("error getting user by UID: %v", err)
	}
	return user, nil
}

// ... (Rest of the code) ...
