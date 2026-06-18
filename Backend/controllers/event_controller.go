package controllers

import (
	"backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetEvents maneja la obtención del catálogo con filtros opcionales (?category=...)
func GetEvents(c *gin.Context) {
	// Lee el parámetro de query opcional desde la URL
	category := c.Query("category")

	events, err := services.GetAllEvents(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los eventos: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetEventByID maneja el detalle de un evento por ID
func GetEventByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de evento inválido"})
		return
	}

	event, err := services.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}