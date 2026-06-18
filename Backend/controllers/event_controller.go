package controllers

import (
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	events, err := services.GetAllEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los eventos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func GetEventByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	event, err := services.GetEventByID(uint(id))
	if err != nil {
		if err.Error() == "evento no encontrado" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}
