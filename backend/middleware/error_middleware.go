package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles errors.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there were any errors
		err := c.Errors.Last()
		if err == nil {
			return
		}

		// Log the error
		log.Printf("Error: %v\n", err.Error())

		// Return an error response
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"message": err.Error(), // Include the error message for debugging purposes
		})
	}
}
