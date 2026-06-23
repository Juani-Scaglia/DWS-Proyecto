package controllers

import (
	"net/http"
	"strconv"

	"backend/services"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

func GetOccupationReportAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID de evento inválido")
		return
	}

	report, err := services.GetOccupationReport(uint(id))
	if err != nil {
		if err.Error() == "evento no encontrado" {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, report)
}
