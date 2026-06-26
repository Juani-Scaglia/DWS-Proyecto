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
	dao.InitDB()

	r := gin.Default()
	r.SetTrustedProxies(nil)

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

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.GET("/events", controllers.GetEvents)
		api.GET(routeEventByID, controllers.GetEventByID)
		api.GET("/events/:id/seats", controllers.GetEventSeats)
		api.GET("/venues", controllers.GetVenues)
		api.GET(routeVenueByID, controllers.GetVenueByID)
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", controllers.LoginUser)

		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/tickets/purchase", controllers.PurchaseTicket)
			protected.GET("/tickets/my-tickets", controllers.GetMyTickets)
			protected.POST("/tickets/:id/cancel", controllers.CancelTicket)
			protected.POST("/tickets/:id/transfer", controllers.TransferTicket)
		}

		admin := api.Group("/admin")
		admin.Use(middlewares.AuthMiddleware())
		admin.Use(middlewares.AdminMiddleware())
		{
			admin.POST("/upload", controllers.UploadImage)
			admin.POST("/events", controllers.CreateEventAdmin)
			admin.PUT(routeEventByID, controllers.UpdateEventAdmin)
			admin.DELETE(routeEventByID, controllers.DeleteEventAdmin)
			admin.POST("/venues", controllers.CreateVenueAdmin)
			admin.PUT(routeVenueByID, controllers.UpdateVenueAdmin)
			admin.DELETE(routeVenueByID, controllers.DeleteVenueAdmin)
			admin.GET("/events/:id/report", controllers.GetOccupationReportAdmin)
		}
	}

	r.Run(":8080")
}
