package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model
	Titulo         string    `json:"titulo" gorm:"type:varchar(150);not null"`
	Descripcion    string    `json:"descripcion" gorm:"type:text"`
	Fecha          time.Time `json:"fecha" gorm:"not null"`
	Categoria      string    `json:"categoria" gorm:"type:varchar(50);not null"`
	Imagen         string    `json:"imagen" gorm:"type:varchar(255)"` // 📸 Foto reglamentaria
	CupoDisponible int       `json:"cupo_disponible" gorm:"not null"`
	VenueID        uint      `json:"venue_id" gorm:"not null"` // 🏛️ Relación con Venue
	Venue          *Venue    `json:"venue,omitempty" gorm:"association_autoupdate:false;association_autocreate:false"`
	Seats          []Seat    `json:"seats,omitempty" gorm:"foreignkey:EventID"`
}
