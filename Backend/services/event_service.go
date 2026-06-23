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
	Precio      float64   `json:"precio" binding:"required,gt=0"`
	Imagen      string    `json:"imagen"`
	VenueID     uint      `json:"venue_id" binding:"required"`
}

func GetAllEvents(category string) ([]domain.Event, error) {
	return dao.GetAllEvents(category)
}

func GetEventByID(id uint) (domain.Event, error) {
	return dao.GetEventByID(id)
}

func CreateEvent(input EventInput) (*domain.Event, error) {
	venue, err := dao.GetVenueByID(input.VenueID)
	if err != nil {
		return nil, err
	}

	var event *domain.Event
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		e := &domain.Event{
			Titulo:      input.Titulo,
			Descripcion: input.Descripcion,
			Categoria:   input.Categoria,
			Fecha:       input.Fecha,
			Lugar:       venue.Nombre + " - " + venue.Direccion,
			Precio:      input.Precio,
			Imagen:      input.Imagen,
			CupoMaximo:  venue.Capacidad,
			CupoDispon:  venue.Capacidad,
			VenueID:     input.VenueID,
		}
		if err := tx.Create(e).Error; err != nil {
			return err
		}
		if err := dao.CreateSeatsForEvent(tx, e.ID, venue.Filas, venue.ColumnasPorFila); err != nil {
			return err
		}
		event = e
		return nil
	})
	return event, err
}

func UpdateEvent(id uint, input EventInput) (*domain.Event, error) {
	if _, err := dao.GetEventByID(id); err != nil {
		return nil, err
	}
	venue, err := dao.GetVenueByID(input.VenueID)
	if err != nil {
		return nil, err
	}
	fields := map[string]interface{}{
		"titulo":      input.Titulo,
		"descripcion": input.Descripcion,
		"categoria":   input.Categoria,
		"fecha":       input.Fecha,
		"lugar":       venue.Nombre + " - " + venue.Direccion,
		"precio":      input.Precio,
		"imagen":      input.Imagen,
		"venue_id":    input.VenueID,
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
		if err := tx.Where("event_id = ?", id).Delete(&domain.Seat{}).Error; err != nil {
			return err
		}
		return tx.Delete(&domain.Event{}, id).Error
	})
}

func GetSeatsByEventID(eventID uint) ([]domain.Seat, error) {
	return dao.GetSeatsByEventID(eventID)
}
