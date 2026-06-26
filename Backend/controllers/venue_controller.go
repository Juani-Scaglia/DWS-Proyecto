package controllers

import (
	"net/http"
	"strconv"

	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

const errVenueIDInvalido = "ID de establecimiento inválido"

// --- RUTAS PÚBLICAS ---

func GetVenues(c *gin.Context) {
	venues, err := services.GetAllVenues()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "No se pudieron obtener los establecimientos: "+err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venues)
}

func GetVenueByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errVenueIDInvalido)
		return
	}
	venue, err := services.GetVenueByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venue)
}

// --- RUTAS DE ADMINISTRACIÓN (ABM) ---

func CreateVenueAdmin(c *gin.Context) {
	var input services.VenueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	venue, err := services.CreateVenue(input)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, venue)
}

func UpdateVenueAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errVenueIDInvalido)
		return
	}
	var input services.VenueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	venue, err := services.UpdateVenue(uint(id), input)
	if err != nil {
		if err.Error() == "establecimiento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venue)
}

func DeleteVenueAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, errVenueIDInvalido)
		return
	}
	if err := services.DeleteVenue(uint(id)); err != nil {
		if err.Error() == "establecimiento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "establecimiento eliminado"})
}
