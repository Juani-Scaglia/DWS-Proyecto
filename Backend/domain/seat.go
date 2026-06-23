package domain

import "github.com/jinzhu/gorm"

type Seat struct {
	gorm.Model
	Numero    string `json:"numero" gorm:"type:varchar(10);not null"`    // Ej: "A-12"
	SeatsioID string `json:"seatsio_id" gorm:"type:varchar(100);unique"` // 🗺️ Conexión al mapa de Juan
	EventID   uint   `json:"event_id" gorm:"not null"`
	TicketID  *uint  `json:"ticket_id" gorm:"default:null"` // Null = Libre
}
