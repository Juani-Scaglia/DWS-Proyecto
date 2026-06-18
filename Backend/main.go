package main

import (
	"backend/controllers" // Reemplazá por el path real de tu proyecto
	"backend/dao"
	"backend/middlewares" // Acá Juan va a meter el validador de JWT

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inicializar la Base de Datos (GORM)
	dao.InitDB()

	// 2. Crear el servidor/enrutador de Gin
	r := gin.Default()

	// 3. Configurar Middleware de CORS (¡Muy importante para que Simon no tenga errores desde React!)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // URL de Vite
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
		// --- RUTAS PÚBLICAS ---

		// Autenticación
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		// Eventos (Catálogo y Detalle) [cite: 44, 48]
		api.GET("/events", controllers.GetEvents)        // Funcionalidad A (Filtros) [cite: 41, 43]
		api.GET("/events/:id", controllers.GetEventByID) // Funcionalidad B [cite: 45]

		// --- RUTAS PROTEGIDAS (Requieren que Juan valide el Token JWT) ---
		// Usamos un grupo que use el middleware de autenticación [cite: 52, 55, 63, 66]
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // <--- Middleware que va a programar Juan
		{
			// Tickets / Entradas
			protected.POST("/tickets/purchase", controllers.PurchaseTicket)     // Funcionalidad C [cite: 49]
			protected.GET("/tickets/my-tickets", controllers.GetMyTickets)      // Funcionalidad D [cite: 53]
			protected.POST("/tickets/:id/cancel", controllers.CancelTicket)     // Funcionalidad E [cite: 61]
			protected.POST("/tickets/:id/transfer", controllers.TransferTicket) // Funcionalidad F [cite: 63]
		}
	}

	// 5. Correr el servidor en el puerto 8080
	r.Run(":8080")
}
