package response

import "github.com/gin-gonic/gin"

// Success sends a successful response with data
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   message,
	})
}

// List sends a paginated list response
func List(c *gin.Context, code int, message string, data interface{}, total int64, page int, limit int) {
	c.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
