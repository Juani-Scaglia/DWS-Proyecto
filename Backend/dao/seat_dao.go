package dao

import (
	"backend/domain/models"
	"errors"

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
	if DB == nil {
		return nil, errors.New(errDBNula)
	}
	var seats []models.Seat
	err := DB.Where("event_id = ?", eventID).Order("fila ASC, numero ASC").Find(&seats).Error
	return seats, err
}

func GetSeatByID(id uint) (models.Seat, error) {
	if DB == nil {
		return models.Seat{}, errors.New(errDBNula)
	}
	var seat models.Seat
	err := DB.First(&seat, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return seat, errors.New("asiento no encontrado")
	}
	return seat, err
}

func OccupySeats(tx *gorm.DB, seatIDs []uint) error {
	return tx.Model(&models.Seat{}).Where("id IN ?", seatIDs).Update("ocupado", true).Error
}

func FreeSeats(tx *gorm.DB, seatIDs []uint) error {
	return tx.Model(&models.Seat{}).Where("id IN ?", seatIDs).Update("ocupado", false).Error
}

func FreeSeatByTicketSeatID(tx *gorm.DB, seatID uint) error {
	return tx.Model(&models.Seat{}).Where("id = ?", seatID).Update("ocupado", false).Error
}
