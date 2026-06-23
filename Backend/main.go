package main

import (
	"backend/controllers"
	"backend/dao"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

const routeEventByID = "/events/:id"
const routeVenueByID = "/venues/:id"

func main() {
	// 1. Inicializar la Base de Datos (GORM)
	dao.InitDB()

	// 2. Crear el servidor/enrutador de Gin
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// 3. Configurar Middleware de CORS
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			origin = "http://localhost:5173"
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ==========================================
	// 4. DEFINICIÓN DE ENDPOINTS
	// ==========================================
	api := r.Group("/api")
	{
		// --- RUTAS PÚBLICAS (Eventos, Venues y Auth) ---
		api.GET("/events", controllers.GetEvents)
		api.GET(routeEventByID, controllers.GetEventByID)
		api.GET("/events/:id/seats", controllers.GetEventSeats)
		api.GET("/venues", controllers.GetVenues)
		api.GET(routeVenueByID, controllers.GetVenueByID)
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		// --- RUTAS PROTEGIDAS (Requieren Token JWT) ---
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/tickets/purchase", controllers.PurchaseTicket)
			protected.GET("/tickets/my-tickets", controllers.GetMyTickets)
			protected.POST("/tickets/:id/cancel", controllers.CancelTicket)
			protected.POST("/tickets/:id/transfer", controllers.TransferTicket)
		}

		// --- RUTAS ADMIN (Requieren Token JWT + rol admin) ---
		admin := api.Group("/admin")
		admin.Use(middlewares.AuthMiddleware())
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.POST("/events", controllers.CreateEventAdmin)
			admin.PUT(routeEventByID, controllers.UpdateEventAdmin)
			admin.DELETE(routeEventByID, controllers.DeleteEventAdmin)
			admin.POST("/venues", controllers.CreateVenueAdmin)
			admin.PUT(routeVenueByID, controllers.UpdateVenueAdmin)
			admin.DELETE(routeVenueByID, controllers.DeleteVenueAdmin)

			// 📊 Reporte de Ocupación para el Administrador (Fase 3)
			admin.GET("/events/:id/report", controllers.GetOccupationReportAdmin)
		}
	}

	// 5. Correr el servidor en el puerto 8080
	r.Run(":8080")
}
