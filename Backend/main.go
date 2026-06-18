package main

import (
<<<<<<< HEAD
	"log"

	"backend/controllers"
	"backend/dao"
	"backend/middlewares"
=======
	"backend/controllers"
	"backend/dao"

	// "backend/middlewares" // Comentado temporalmente hasta que Juan lo cree
>>>>>>> feature/eventos

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Inicializar la Base de Datos (GORM)
	dao.InitDB()

	// 2. Crear el servidor/enrutador de Gin
	r := gin.Default()

<<<<<<< HEAD
	// 3. Configurar Middleware de CORS para conectar con el React de Simon
=======
	// 3. Configurar Middleware de CORS
>>>>>>> feature/eventos
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
<<<<<<< HEAD
		// --- RUTAS PÚBLICAS ---
		
		// Autenticación [cite: 38]
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		// Eventos (Catálogo y Detalle) [cite: 44, 48]
		api.GET("/events", controllers.GetEvents)       // Funcionalidad A [cite: 41]
		api.GET("/events/:id", controllers.GetEventByID) // Funcionalidad B [cite: 45]

		// --- RUTAS PROTEGIDAS (Requieren Token JWT) ---
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware()) // Middleware de Juan [cite: 194]
		{
			// Tickets / Entradas [cite: 27]
			protected.POST("/tickets/purchase", controllers.PurchaseTicket)     // Funcionalidad C [cite: 49]
			protected.GET("/tickets/my-tickets", controllers.GetMyTickets)      // Funcionalidad D [cite: 53]
			protected.POST("/tickets/:id/cancel", controllers.CancelTicket)     // Funcionalidad E [cite: 61]
			protected.POST("/tickets/:id/transfer", controllers.TransferTicket) // Funcionalidad F [cite: 63]
		}
=======
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
>>>>>>> feature/eventos
	}

	// 5. Correr el servidor en el puerto 8080
	log.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
