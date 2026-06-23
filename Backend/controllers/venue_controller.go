package controllers

import (
	"net/http"
	"strconv"

	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func GetVenues(c *gin.Context) {
	venues, err := services.GetAllVenues()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venues)
}

func GetVenueByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de establecimiento inválido")
		return
	}
	venue, err := services.GetVenueByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venue)
}

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
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de establecimiento inválido")
		return
	}
	var input services.VenueInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	venue, err := services.UpdateVenue(uint(id), input)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "establecimiento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, errMsg)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, errMsg)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, venue)
}

func DeleteVenueAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de establecimiento inválido")
		return
	}
	if err := services.DeleteVenue(uint(id)); err != nil {
		errMsg := err.Error()
		if errMsg == "establecimiento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, errMsg)
			return
		}
		utils.ErrorResponse(c, http.StatusBadRequest, errMsg)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "establecimiento eliminado"})
}
