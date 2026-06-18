package main

import (
	"log"

	"backend/controllers"
	"backend/dao"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	dao.InitDB()

	r := gin.Default()

	// CORS para conectar con el frontend React
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

	api := r.Group("/api")
	{
		// --- RUTAS PÚBLICAS ---
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		api.GET("/events", controllers.GetEvents)
		api.GET("/events/:id", controllers.GetEventByID)

		// --- RUTAS PROTEGIDAS ---
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/tickets/purchase", controllers.PurchaseTicket)
			protected.GET("/tickets/my-tickets", controllers.GetMyTickets)
			protected.POST("/tickets/:id/cancel", controllers.CancelTicket)
			protected.POST("/tickets/:id/transfer", controllers.TransferTicket)
		}
	}

	log.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
