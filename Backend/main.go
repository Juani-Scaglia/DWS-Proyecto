package main

import (
	"backend/controllers"
	"backend/dao"

	// "backend/middlewares" // Comentado temporalmente hasta que Juan lo cree

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inicializar la Base de Datos (GORM)
	dao.InitDB()

	// 2. Crear el servidor/enrutador de Gin
	r := gin.Default()

	// 3. Configurar Middleware de CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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
		// --- RUTAS PÚBLICAS DE EVENTOS (¡TUS ENDPOINTS!) ---
		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventByID)

		// --- COMENTADO TEMPORALMENTE PARA QUE PUEDAS COMPILAR Y PROBAR ---
		/*
			// Autenticación
			api.POST("/auth/register", controllers.RegisterUser)
			api.POST("/auth/login", controllers.LoginUser)

			// Rutas protegidas
			protected := api.Group("/")
			protected.Use(middlewares.AuthMiddleware())
			{
				protected.POST("/tickets/purchase", controllers.PurchaseTicket)
				protected.GET("/tickets/my-tickets", controllers.GetMyTickets)
				protected.POST("/tickets/:id/cancel", controllers.CancelTicket)
				protected.POST("/tickets/:id/transfer", controllers.TransferTicket)
			}
		*/
	}

	// 5. Correr el servidor en el puerto 8080
	r.Run(":8080")
}
