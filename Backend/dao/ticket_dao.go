package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func CreateTicket(tx *gorm.DB, ticket *models.Ticket) error {
	return tx.Create(ticket).Error
}

func DecrementCupo(tx *gorm.DB, eventID uint) error {
	result := tx.Model(&models.Event{}).
		Where("id = ? AND cupo_disponible > 0", eventID).
		UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible - 1"))
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
	err := DB.Preload("Event").Where("user_id = ?", userID).Find(&tickets).Error
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
