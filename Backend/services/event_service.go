package services

import (
	"backend/dao"
	"backend/domain/models"
)

func GetAllEvents() ([]domain.Event, error) {
	return dao.GetAllEvents()
}

func GetEventByID(id uint) (domain.Event, error) {
	return dao.GetEventByID(id)
}