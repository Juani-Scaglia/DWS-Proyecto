package domain

import "github.com/jinzhu/gorm"

type Venue struct {
	gorm.Model
	Nombre      string  `json:"nombre" gorm:"type:varchar(100);not null"`
	Dimensiones string  `json:"dimensiones" gorm:"type:varchar(50)"`
	Capacidad   int     `json:"capacidad" gorm:"not null"`
	Events      []Event `json:"events,omitempty" gorm:"foreignkey:VenueID"`
}
