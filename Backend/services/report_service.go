package services

import (
	"backend/dao"
	domain "backend/domain/models"
	"errors"
)

type OccupationReport struct {
	EventID           uint    `json:"event_id"`
	EventName         string  `json:"event_name"`
	VenueName         string  `json:"venue_name"`
	TicketsIssued     int64   `json:"tickets_issued"`
	TotalCapacity     int64   `json:"total_capacity"`
	OccupationPercent float64 `json:"occupation_percent"`
}

// GetOccupationReport calcula el ratio de entradas emitidas vs capacidad para un evento
func GetOccupationReport(eventID uint) (*OccupationReport, error) {
	if dao.DB == nil {
		return nil, errors.New("base de datos no inicializada")
	}

	// 1. Buscar el evento
	var event domain.Event
	if err := dao.DB.First(&event, eventID).Error; err != nil {
		return nil, errors.New("evento no encontrado")
	}

	// 2. Buscar el establecimiento (Venue) usando VenueID
	var venue domain.Venue
	if err := dao.DB.First(&venue, event.VenueID).Error; err != nil {
		return nil, errors.New("establecimiento asociado no encontrado")
	}

	// 3. Contar cuántos tickets activos se emitieron para este evento
	var ticketsIssued int64
	if err := dao.DB.Model(&domain.Ticket{}).
		Where("event_id = ? AND estado = 'activo'", eventID).
		Count(&ticketsIssued).Error; err != nil {
		return nil, err
	}

	// 4. Calcular el porcentaje de ocupación evadiendo la división por cero
	var percent float64
	if venue.Capacidad > 0 {
		percent = (float64(ticketsIssued) / float64(venue.Capacidad)) * 100
	}

	return &OccupationReport{
		EventID:           event.ID,
		EventName:         event.Titulo, // Usamos Titulo que es el campo real de tu modelo
		VenueName:         venue.Nombre,
		TicketsIssued:     ticketsIssued,
		TotalCapacity:     int64(venue.Capacidad),
		OccupationPercent: percent,
	}, nil
}
