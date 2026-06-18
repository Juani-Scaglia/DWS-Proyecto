package controllers

import (
	"net/http"
	"strconv"

	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

const errEventoIDInvalido = "ID de evento inválido"

func GetEvents(c *gin.Context) {
	category := c.Query("category")
	events, err := services.GetAllEvents(category)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudieron obtener los eventos: "+err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errEventoIDInvalido)
		return
	}
	event, err := services.GetEventByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, event)
}

func CreateEventAdmin(c *gin.Context) {
	var input services.EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	event, err := services.CreateEvent(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, event)
}

func UpdateEventAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errEventoIDInvalido)
		return
	}
	var input services.EventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	event, err := services.UpdateEvent(uint(id), input)
	if err != nil {
		if err.Error() == "evento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, event)
}

func DeleteEventAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errEventoIDInvalido)
		return
	}
	if err := services.DeleteEvent(uint(id)); err != nil {
		if err.Error() == "evento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "evento eliminado"})
}
