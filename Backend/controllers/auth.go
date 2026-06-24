package controllers

import (
	"net/http"
	"strings"

	"backend/services"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var input services.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		msg := err.Error()
		if strings.Contains(msg, "'Password'") && strings.Contains(msg, "'min'") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "La contraseña debe tener al menos 6 caracteres"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	user, err := services.Register(input)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario registrado con éxito",
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"nombre":   user.Nombre,
			"apellido": user.Apellido,
			"rol":      user.Rol,
		},
	})
}

func LoginUser(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := services.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"nombre":   user.Nombre,
			"apellido": user.Apellido,
			"dni":      user.DNI,
			"rol":      user.Rol,
		},
	})
}
