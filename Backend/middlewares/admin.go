package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, _ := c.Get("rol")
		if rol != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acceso restringido a administradores"})
			return
		}
		c.Next()
	}
}
