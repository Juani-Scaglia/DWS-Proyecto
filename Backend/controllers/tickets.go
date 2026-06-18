package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PurchaseTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "purchase ticket - TODO"})
}

func GetMyTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get my tickets - TODO"})
}

func CancelTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "cancel ticket - TODO"})
}

func TransferTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "transfer ticket - TODO"})
}
