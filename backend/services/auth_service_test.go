package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	userID := "testuser"
	token, err := GenerateJWT(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

// Add more tests for other auth service functions
