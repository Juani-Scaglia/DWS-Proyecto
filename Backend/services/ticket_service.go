package services

import (
	"errors"
	"fmt"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"gorm.io/gorm"
)

func PurchaseTickets(userID uint, eventID uint, seatIDs []uint) ([]domain.Ticket, error) {
	if len(seatIDs) == 0 {
		return nil, errors.New("debés seleccionar al menos un asiento")
	}

	userCount, err := dao.GetActiveTicketCountByEventAndUser(userID, eventID)
	if err != nil {
		return nil, err
	}
	if int(userCount)+len(seatIDs) > 10 {
		return nil, fmt.Errorf("límite de 10 entradas por evento (tenés %d)", userCount)
	}

	for _, sid := range seatIDs {
		seat, err := dao.GetSeatByID(sid)
		if err != nil {
			return nil, err
		}
		if seat.EventID != eventID {
			return nil, errors.New("el asiento no pertenece a este evento")
		}
		if seat.Ocupado {
			return nil, fmt.Errorf("el asiento %s%d ya está ocupado", seat.Fila, seat.Numero)
		}
	}

	var tickets []domain.Ticket
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := dao.DecrementCupo(tx, eventID, len(seatIDs)); err != nil {
			return err
		}
		if err := dao.OccupySeats(tx, seatIDs); err != nil {
			return err
		}
		for _, sid := range seatIDs {
			seatID := sid
			t := domain.Ticket{
				UserID:    userID,
				EventID:   eventID,
				SeatID:    &seatID,
				FechaComp: time.Now(),
				Estado:    "activo",
			}
			if err := tx.Create(&t).Error; err != nil {
				return err
			}
			tickets = append(tickets, t)
		}
		return nil
	})
	return tickets, err
}

func GetMyTickets(userID uint) ([]domain.Ticket, error) {
	return dao.GetTicketsByUserID(userID)
}

func CancelTicket(userID uint, ticketID uint) error {
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
		if err := dao.IncrementCupo(tx, ticket.EventID); err != nil {
			return err
		}
		if ticket.SeatID != nil {
			return dao.FreeSeat(tx, *ticket.SeatID)
		}
		return nil
	})
}

func TransferTicket(ownerID uint, ticketID uint, targetDNI string) error {
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

func GetEventReport(eventID uint) (map[string]interface{}, error) {
	return dao.GetEventReport(eventID)
}
