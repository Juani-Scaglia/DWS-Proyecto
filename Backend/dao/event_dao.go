package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

const errDBNula = "base de datos no inicializada"

func GetAllEvents(category string) ([]models.Event, error) {
	if DB == nil {
		return nil, errors.New(errDBNula)
	}
	var events []models.Event

	if category != "" {
		result := DB.Where("categoria = ?", category).Find(&events)
		if result.Error != nil {
			return nil, result.Error
		}
		return events, nil
	}

	result := DB.Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func GetEventByID(id uint) (models.Event, error) {
	if DB == nil {
		return models.Event{}, errors.New(errDBNula)
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

func CreateEvent(event *models.Event) error {
	if DB == nil {
		return errors.New(errDBNula)
	}
	return DB.Create(event).Error
}

func UpdateEvent(id uint, fields map[string]interface{}) error {
	return DB.Model(&models.Event{}).Where("id = ?", id).Updates(fields).Error
}

func DeleteEvent(id uint) error {
	return DB.Delete(&models.Event{}, id).Error
}
