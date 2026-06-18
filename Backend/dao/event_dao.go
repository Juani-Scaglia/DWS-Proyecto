package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func GetAllEvents() ([]domain.Event, error) {
	var events []domain.Event
	result := DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func GetEventByID(id uint) (domain.Event, error) {
	var event domain.Event
	result := DB.First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return event, errors.New("evento no encontrado")
		}
		return event, result.Error
	}
	return event, nil
}
