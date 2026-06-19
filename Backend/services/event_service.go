package services

import (
	"time"

	"backend/dao"
	domain "backend/domain/models"
	"gorm.io/gorm"
)

type EventInput struct {
	Titulo      string    `json:"titulo" binding:"required"`
	Descripcion string    `json:"descripcion"`
	Categoria   string    `json:"categoria" binding:"required"`
	Fecha       time.Time `json:"fecha" binding:"required"`
	Lugar       string    `json:"lugar" binding:"required"`
	Precio      float64   `json:"precio" binding:"required,gt=0"`
	CupoMaximo  int       `json:"cupo_maximo" binding:"required,gt=0"`
}

func GetAllEvents(category string) ([]domain.Event, error) {
	return dao.GetAllEvents(category)
}

func GetEventByID(id uint) (domain.Event, error) {
	return dao.GetEventByID(id)
}

func CreateEvent(input EventInput) (*domain.Event, error) {
	event := &domain.Event{
		Titulo:      input.Titulo,
		Descripcion: input.Descripcion,
		Categoria:   input.Categoria,
		Fecha:       input.Fecha,
		Lugar:       input.Lugar,
		Precio:      input.Precio,
		CupoMaximo:  input.CupoMaximo,
		CupoDispon:  input.CupoMaximo,
	}
	if err := dao.CreateEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}

func UpdateEvent(id uint, input EventInput) (*domain.Event, error) {
	if _, err := dao.GetEventByID(id); err != nil {
		return nil, err
	}
	fields := map[string]interface{}{
		"titulo":      input.Titulo,
		"descripcion": input.Descripcion,
		"categoria":   input.Categoria,
		"fecha":       input.Fecha,
		"lugar":       input.Lugar,
		"precio":      input.Precio,
		"cupo_maximo": input.CupoMaximo,
	}
	if err := dao.UpdateEvent(id, fields); err != nil {
		return nil, err
	}
	event, err := dao.GetEventByID(id)
	return &event, err
}

func DeleteEvent(id uint) error {
	if _, err := dao.GetEventByID(id); err != nil {
		return err
	}
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("event_id = ?", id).Delete(&domain.Ticket{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Event{}, id).Error
	})
}
