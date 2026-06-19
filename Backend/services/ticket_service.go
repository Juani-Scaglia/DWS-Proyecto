package services

import (
	"errors"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"gorm.io/gorm"
)

func PurchaseTicket(userID uint, eventID uint) (*domain.Ticket, error) {
	if dao.DB == nil {
		return nil, errors.New("base de datos no inicializada")
	}

	var userCount int64
	if err := dao.DB.Model(&domain.Ticket{}).
		Where("user_id = ? AND event_id = ? AND estado = 'activo'", userID, eventID).
		Count(&userCount).Error; err != nil {
		return nil, err
	}
	if userCount >= 10 {
		return nil, errors.New("alcanzaste el límite de 10 entradas para este evento")
	}

	var ticket *domain.Ticket

	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := dao.DecrementCupo(tx, eventID); err != nil {
			return err
		}
		t := &domain.Ticket{
			UserID:    userID,
			EventID:   eventID,
			FechaComp: time.Now(),
			Estado:    "activo",
		}
		if err := dao.CreateTicket(tx, t); err != nil {
			return err
		}
		ticket = t
		return nil
	})

	return ticket, err
}

func GetMyTickets(userID uint) ([]domain.Ticket, error) {
	return dao.GetTicketsByUserID(userID)
}

func CancelTicket(userID uint, ticketID uint) error {
	if dao.DB == nil {
		return errors.New("base de datos no inicializada")
	}
	ticket, err := dao.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.UserID != userID {
		return errors.New("no autorizado")
	}
	if ticket.Estado != "activo" {
		return errors.New("el ticket no está activo")
	}

	return dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := dao.UpdateTicketEstado(tx, ticketID, "cancelado"); err != nil {
			return err
		}
		return dao.IncrementCupo(tx, ticket.EventID)
	})
}

func TransferTicket(ownerID uint, ticketID uint, targetDNI string) error {
	if dao.DB == nil {
		return errors.New("base de datos no inicializada")
	}
	ticket, err := dao.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.UserID != ownerID {
		return errors.New("no autorizado")
	}
	if ticket.Estado != "activo" {
		return errors.New("el ticket no está activo")
	}

	targetUser, err := dao.GetUserByDNI(targetDNI)
	if err != nil {
		return err
	}
	if targetUser.ID == ownerID {
		return errors.New("no podés transferirte el ticket a vos mismo")
	}

	return dao.DB.Transaction(func(tx *gorm.DB) error {
		return dao.TransferTicketOwner(tx, ticketID, targetUser.ID)
	})
}
