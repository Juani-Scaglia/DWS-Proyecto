package controllers

import (
	"net/http"
	"strconv"

	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type PurchaseInput struct {
	EventID uint   `json:"event_id" binding:"required"`
	SeatIDs []uint `json:"seat_ids" binding:"required,min=1"`
}

type TransferInput struct {
	DNI string `json:"dni" binding:"required"`
}

func PurchaseTicket(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tickets, err := services.PurchaseTickets(userID, input.EventID, input.SeatIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, tickets)
}

func GetMyTickets(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	tickets, err := services.GetMyTickets(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, tickets)
}

func CancelTicket(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de ticket inválido")
		return
	}

	if err := services.CancelTicket(userID, uint(ticketID)); err != nil {
		switch err.Error() {
		case "no autorizado":
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		case "ticket no encontrado":
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		default:
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "ticket cancelado"})
}

func TransferTicket(c *gin.Context) {
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de ticket inválido")
		return
	}

	var input TransferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.TransferTicket(userID, uint(ticketID), input.DNI); err != nil {
		switch err.Error() {
		case "no autorizado":
			utils.ErrorResponse(c, http.StatusForbidden, err.Error())
		case "ticket no encontrado", "usuario con ese DNI no encontrado":
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		default:
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "ticket transferido"})
}

func GetEventReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID inválido")
		return
	}
	report, err := services.GetEventReport(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, report)
}
