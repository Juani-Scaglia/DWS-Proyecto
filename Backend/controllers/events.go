package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEvents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get events - TODO"})
}

func GetEventByID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get event by id - TODO"})
}
