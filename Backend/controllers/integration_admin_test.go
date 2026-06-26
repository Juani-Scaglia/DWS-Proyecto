package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func TestAdminRoutes_SecurityAndPermissions(t *testing.T) {
	// 1. Configurar Gin en modo test
	gin.SetMode(gin.TestMode)

	t.Run("Crear evento sin Token - Debe dar 401 Unauthorized", func(t *testing.T) {
		r := gin.Default()
		r.POST("/api/admin/events", middlewares.AuthMiddleware(), controllers.CreateEventAdmin)

		req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(`{"titulo":"Concierto Infiltrado"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Se esperaba código 401, pero se obtuvo %d", w.Code)
		}
	})

	t.Run("Cliente intenta crear evento - Debe dar 403 Forbidden", func(t *testing.T) {
		r := gin.Default()
		r.POST("/api/admin/events", func(c *gin.Context) {
			// Simulamos un token válido inyectado de rol Cliente
			c.Set("user_id", uint(1))
			c.Set("rol", "Cliente")
			c.Next()
		}, middlewares.AdminMiddleware(), controllers.CreateEventAdmin)

		req, _ := http.NewRequest("POST", "/api/admin/events", strings.NewReader(`{"titulo":"Concierto Prohibido"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Se esperaba código 403 (Forbidden) para un Cliente, pero se obtuvo %d", w.Code)
		}
	})

	t.Run("Cliente intenta ver reporte de métricas - Debe dar 403 Forbidden", func(t *testing.T) {
		r := gin.Default()
		r.GET("/api/admin/events/:id/report", func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("rol", "Cliente")
			c.Next()
		}, middlewares.AdminMiddleware(), controllers.GetOccupationReportAdmin)

		req, _ := http.NewRequest("GET", "/api/admin/events/1/report", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Se esperaba código 403 al intentar acceder al reporte siendo Cliente, pero se obtuvo %d", w.Code)
		}
	})
}
