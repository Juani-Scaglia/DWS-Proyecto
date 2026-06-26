package dao

import (
	"backend/domain/models"
	"errors"
	"math"

	"gorm.io/gorm"
)

type SectorDef struct {
	Nombre    string
	Capacidad int
}

func CreateSeatsForEvent(tx *gorm.DB, eventID uint, sectores []SectorDef) error {
	batch := make([]models.Seat, 0, 500)
	for _, sec := range sectores {
		if sec.Capacidad <= 0 {
			continue
		}
		cols := int(math.Ceil(math.Sqrt(float64(sec.Capacidad))))
		if cols > 50 {
			cols = 50
		}
		filas := int(math.Ceil(float64(sec.Capacidad) / float64(cols)))
		count := 0
		for f := 0; f < filas && count < sec.Capacidad; f++ {
			fila := seatRowLabel(f)
			for c := 1; c <= cols && count < sec.Capacidad; c++ {
				batch = append(batch, models.Seat{
					EventID: eventID,
					Sector:  sec.Nombre,
					Fila:    fila,
					Numero:  c,
					Ocupado: false,
				})
				count++
				if len(batch) >= 500 {
					if err := tx.Create(&batch).Error; err != nil {
						return err
					}
					batch = batch[:0]
				}
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
	if DB == nil {
		return nil, errors.New(errDBNula)
	}
	var seats []models.Seat
	err := DB.Where("event_id = ?", eventID).Order("sector ASC, fila ASC, numero ASC").Find(&seats).Error
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
