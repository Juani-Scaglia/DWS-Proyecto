package dao

import (
	"backend/domain/models"

	"gorm.io/gorm"
)

func CreateSeatsForEvent(tx *gorm.DB, eventID uint, filas int, columnas int) error {
	batch := make([]models.Seat, 0, 500)
	for f := 0; f < filas; f++ {
		fila := seatRowLabel(f)
		for c := 1; c <= columnas; c++ {
			batch = append(batch, models.Seat{
				EventID: eventID,
				Fila:    fila,
				Numero:  c,
				Ocupado: false,
			})
			if len(batch) >= 500 {
				if err := tx.Create(&batch).Error; err != nil {
					return err
				}
				batch = batch[:0]
			}
		}
	}
	if len(batch) > 0 {
		return tx.Create(&batch).Error
	}
	return nil
}

func seatRowLabel(index int) string {
	label := ""
	for index >= 0 {
		label = string(rune('A'+index%26)) + label
		index = index/26 - 1
	}
	return label
}

func GetSeatsByEventID(eventID uint) ([]models.Seat, error) {
	var seats []models.Seat
	err := DB.Where("event_id = ?", eventID).Order("fila ASC, numero ASC").Find(&seats).Error
	return seats, err
}
