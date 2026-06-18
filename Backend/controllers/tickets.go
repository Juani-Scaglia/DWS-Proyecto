package controllers

import (
	"net/http"
	"strconv"

	"backend/services"

	"github.com/gin-gonic/gin"
)

type PurchaseInput struct {
	EventID uint `json:"event_id" binding:"required"`
}

type TransferInput struct {
	DNI string `json:"dni" binding:"required"`
}

func PurchaseTicket(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := services.PurchaseTicket(userID, input.EventID)
	if err != nil {
		if err.Error() == "sin cupo disponible" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func GetMyTickets(c *gin.Context) {
	userID := c.GetUint("user_id")

	tickets, err := services.GetMyTickets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func CancelTicket(c *gin.Context) {
	userID := c.GetUint("user_id")

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket inválido"})
		return
	}

	if err := services.CancelTicket(userID, uint(ticketID)); err != nil {
		switch err.Error() {
		case "no autorizado":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "ticket no encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket cancelado"})
}

func TransferTicket(c *gin.Context) {
	userID := c.GetUint("user_id")

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket inválido"})
		return
	}

	var input TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.TransferTicket(userID, uint(ticketID), input.DNI); err != nil {
		switch err.Error() {
		case "no autorizado":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "ticket no encontrado", "usuario con ese DNI no encontrado":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket transferido"})
}
