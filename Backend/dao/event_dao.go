package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func GetAllEvents(category string) ([]models.Event, error) {
	if DB == nil {
		return nil, errors.New(errDBNula)
	}
	var events []models.Event
	q := DB.Preload("Venue")
	if category != "" {
		q = q.Where("categoria = ?", category)
	}
	return events, q.Find(&events).Error
}

func GetEventByID(id uint) (models.Event, error) {
	if DB == nil {
		return models.Event{}, errors.New(errDBNula)
	}
	var event models.Event
	err := DB.Preload("Venue").First(&event, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return event, errors.New("evento no encontrado")
	}
	return event, err
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
