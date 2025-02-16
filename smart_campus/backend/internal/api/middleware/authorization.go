package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
