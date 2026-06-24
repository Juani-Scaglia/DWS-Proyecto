package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"backend/utils"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("imagen")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "No se recibio ninguna imagen")
		return
	}

	ext := filepath.Ext(file.Filename)
	nombre := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	dir := "uploads"
	os.MkdirAll(dir, os.ModePerm)

	ruta := filepath.Join(dir, nombre)
	if err := c.SaveUploadedFile(file, ruta); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error al guardar la imagen")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/uploads/%s", nombre)
	utils.SuccessResponse(c, http.StatusOK, gin.H{"url": url})
}
