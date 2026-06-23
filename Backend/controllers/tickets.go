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
	userID := c.GetUint("user_id")

	var input PurchaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tickets, err := services.PurchaseTickets(userID, input.EventID, input.SeatIDs)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "sin cupo disponible" {
			utils.ErrorResponse(c, http.StatusConflict, errMsg)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, errMsg)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, tickets)
}

func GetMyTickets(c *gin.Context) {
	userID := c.GetUint("user_id")

	tickets, err := services.GetMyTickets(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, tickets)
}

func CancelTicket(c *gin.Context) {
	userID := c.GetUint("user_id")

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
	userID := c.GetUint("user_id")

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

func GetEventSeats(c *gin.Context) {
	eventID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de evento inválido")
		return
	}
	seats, err := services.GetSeatsByEventID(uint(eventID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, seats)
}
