package services

import (
	"backend/dao"
	"backend/domain/models"
)

// GetAllEvents pasa el filtro de categoría hacia el DAO
func GetAllEvents(category string) ([]models.Event, error) {
	return dao.GetAllEvents(category)
}

func GetEventByID(id uint) (models.Event, error) {
	return dao.GetEventByID(id)
}