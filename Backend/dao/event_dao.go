package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

// GetAllEvents trae los eventos, opcionalmente filtrados por categoría (Requisito Regularidad)
func GetAllEvents(category string) ([]models.Event, error) {
	if DB == nil {
		return nil, errors.New("base de datos no inicializada")
	}
	var events []models.Event

	if category != "" {
		result := DB.Where("category = ?", category).Find(&events)
		if result.Error != nil {
			return nil, result.Error
		}
		return events, nil
	}

	// Si no manda categoría, trae todos de forma general
	result := DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// GetEventByID busca un evento puntual a partir de su ID único
func GetEventByID(id uint) (models.Event, error) {
	if DB == nil {
		return models.Event{}, errors.New("base de datos no inicializada")
	}
	var event models.Event
	result := DB.First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return event, errors.New("evento no encontrado")
		}
		return event, result.Error
	}
	return event, nil
}
