package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware is a middleware that authenticates requests using JWT.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// Extract the token from the header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Verify the token
		token, err := verifyJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Extract the user ID from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			return
		}

		// Set the user ID in the context for later use
		c.Set("user_id", userID)

		// Continue to the next handler
		c.Next()
	}
}

// verifyJWT verifies the JWT token.
func verifyJWT(tokenString string) (*jwt.Token, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Get the secret key from the environment variables
		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			secretKey = "secret" // Use a default secret key if not set in environment variables
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
