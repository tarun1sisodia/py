package middleware

import (
	"net/http"
	"smart_campus/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// Extract token from Bearer header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			return
		}

		// Validate token
		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		// Set user info in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("userRole")
		if userRole == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		isAllowed := false
		for _, role := range roles {
			if userRole == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			return
		}

		c.Next()
	}
}

func RequireDevice() gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Device ID required",
			})
			return
		}

		userID := c.GetString("userID")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		// TODO: Verify device binding
		// deviceService := c.MustGet("deviceService").(*services.DeviceService)
		// isValid, err := deviceService.VerifyDevice(userID, deviceID)
		// if err != nil {
		//     c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		//         "error": "Error verifying device",
		//     })
		//     return
		// }
		// if !isValid {
		//     c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		//         "error": "Invalid device",
		//     })
		//     return
		// }

		c.Set("deviceID", deviceID)
		c.Next()
	}
}
