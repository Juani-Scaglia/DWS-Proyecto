package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, exists := c.Get("rol")
		
		// Valida tanto si viene en minúscula "admin" como "Administrador" con mayúscula
		if !exists || (rol != "admin" && rol != "Administrador") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: se requieren permisos de Administrador",
			})
			return
		}
		
		c.Next()
	}
}