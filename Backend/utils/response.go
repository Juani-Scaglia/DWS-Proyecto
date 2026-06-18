package utils

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
