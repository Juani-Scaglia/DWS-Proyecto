package services

import (
	"errors"
	"fmt"
	"time"

	"backend/dao"
	domain "backend/domain/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PurchaseTickets(userID uint, eventID uint, seatIDs []uint) ([]domain.Ticket, error) {
	if dao.DB == nil {
		return nil, errors.New("base de datos no inicializada")
	}

	var userCount int64
	if err := dao.DB.Model(&domain.Ticket{}).
		Where("user_id = ? AND event_id = ? AND estado = 'activo'", userID, eventID).
		Count(&userCount).Error; err != nil {
		return nil, err
	}
	if int(userCount)+len(seatIDs) > 10 {
		return nil, errors.New("alcanzaste el límite de 10 entradas para este evento")
	}

	var tickets []domain.Ticket

	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		var seats []domain.Seat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id IN ?", seatIDs).
			Find(&seats).Error; err != nil {
			return err
		}

		if len(seats) != len(seatIDs) {
			return errors.New("uno o más IDs de asiento no existen")
		}

		for _, seat := range seats {
			if seat.EventID != eventID {
				return fmt.Errorf("el asiento %s%d no pertenece a este evento", seat.Fila, seat.Numero)
			}
			if seat.Ocupado {
				return fmt.Errorf("el asiento %s%d ya está ocupado", seat.Fila, seat.Numero)
			}
		}

		result := tx.Model(&domain.Event{}).
			Where("id = ? AND cupo_disponible >= ?", eventID, len(seatIDs)).
			UpdateColumn("cupo_disponible", gorm.Expr("cupo_disponible - ?", len(seatIDs)))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("sin cupo disponible para la cantidad solicitada")
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
		if ticket.SeatID != nil {
			if err := dao.FreeSeatByTicketSeatID(tx, *ticket.SeatID); err != nil {
				return err
			}
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
