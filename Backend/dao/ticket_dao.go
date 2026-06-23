package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func CreateTicket(tx *gorm.DB, ticket *models.Ticket) error {
	return tx.Create(ticket).Error
}

func DecrementCupo(tx *gorm.DB, eventID uint, cantidad int) error {
	result := tx.Model(&models.Event{}).
		Where("id = ? AND cupo_disponible >= ?", eventID, cantidad).
		UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible - ?", cantidad))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("sin cupo disponible")
	}
	return nil
}

func IncrementCupo(tx *gorm.DB, eventID uint) error {
	return tx.Model(&models.Event{}).
		Where("id = ?", eventID).
		UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible + 1")).Error
}

func GetTicketsByUserID(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := DB.Preload("Event").Preload("Seat").Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

func GetTicketByID(id uint) (models.Ticket, error) {
	var ticket models.Ticket
	err := DB.First(&ticket, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ticket, errors.New("ticket no encontrado")
	}
	return ticket, err
}

func UpdateTicketEstado(tx *gorm.DB, ticketID uint, estado string) error {
	return tx.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("estado", estado).Error
}

func TransferTicketOwner(tx *gorm.DB, ticketID uint, newUserID uint) error {
	return tx.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("user_id", newUserID).Error
}

func GetUserByDNI(dni string) (models.User, error) {
	var user models.User
	err := DB.Where("dni = ?", dni).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("usuario con ese DNI no encontrado")
	}
	return user, err
}

func OccupySeats(tx *gorm.DB, seatIDs []uint) error {
	return tx.Model(&models.Seat{}).Where("id IN ?", seatIDs).Update("ocupado", true).Error
}

func FreeSeat(tx *gorm.DB, seatID uint) error {
	return tx.Model(&models.Seat{}).Where("id = ?", seatID).Update("ocupado", false).Error
}

func GetSeatByID(id uint) (models.Seat, error) {
	var seat models.Seat
	err := DB.First(&seat, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return seat, errors.New("asiento no encontrado")
	}
	return seat, err
}

func GetActiveTicketCountByEventAndUser(userID, eventID uint) (int64, error) {
	var count int64
	err := DB.Model(&models.Ticket{}).
		Where("user_id = ? AND event_id = ? AND estado = 'activo'", userID, eventID).
		Count(&count).Error
	return count, err
}

func GetEventReport(eventID uint) (map[string]interface{}, error) {
	var event models.Event
	if err := DB.Preload("Venue").First(&event, eventID).Error; err != nil {
		return nil, errors.New("evento no encontrado")
	}

	var vendidos int64
	DB.Model(&models.Ticket{}).Where("event_id = ? AND estado = 'activo'", eventID).Count(&vendidos)

	var cancelados int64
	DB.Model(&models.Ticket{}).Where("event_id = ? AND estado = 'cancelado'", eventID).Count(&cancelados)

	var asientosOcupados int64
	DB.Model(&models.Seat{}).Where("event_id = ? AND ocupado = true", eventID).Count(&asientosOcupados)

	var asientosTotales int64
	DB.Model(&models.Seat{}).Where("event_id = ?", eventID).Count(&asientosTotales)

	porcentaje := 0.0
	if event.CupoMaximo > 0 {
		porcentaje = float64(vendidos) / float64(event.CupoMaximo) * 100
	}

	return map[string]interface{}{
		"evento":             event,
		"entradas_vendidas":  vendidos,
		"entradas_canceladas": cancelados,
		"asientos_ocupados":  asientosOcupados,
		"asientos_totales":   asientosTotales,
		"porcentaje_ocupacion": porcentaje,
	}, nil
}
