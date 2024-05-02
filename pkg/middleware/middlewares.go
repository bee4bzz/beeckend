package middleware

import "github.com/gin-gonic/gin"

func JSONContentType(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}
