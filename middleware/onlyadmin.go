package middleware

import "github.com/gin-gonic/gin"

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role != "admin" {
			c.JSON(403, gin.H{"error": "forbidden: admin only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
