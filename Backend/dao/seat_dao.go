package dao

import (
	"backend/domain/models"

	"gorm.io/gorm"
)

func CreateSeatsForEvent(tx *gorm.DB, eventID uint, filas int, columnas int) error {
	seats := make([]models.Seat, 0, filas*columnas)
	for f := 0; f < filas; f++ {
		fila := string(rune('A' + f))
		for c := 1; c <= columnas; c++ {
			seats = append(seats, models.Seat{
				EventID: eventID,
				Fila:    fila,
				Numero:  c,
				Ocupado: false,
			})
		}
	}
	return tx.Create(&seats).Error
}

func GetSeatsByEventID(eventID uint) ([]models.Seat, error) {
	var seats []models.Seat
	err := DB.Where("event_id = ?", eventID).Order("fila ASC, numero ASC").Find(&seats).Error
	return seats, err
}
