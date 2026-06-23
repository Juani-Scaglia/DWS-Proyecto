package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(191);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Nombre    string    `gorm:"type:varchar(100);not null" json:"nombre"`
	Apellido  string    `gorm:"type:varchar(100);not null" json:"apellido"`
	Rol       string    `gorm:"type:varchar(50);not null;default:'cliente'" json:"rol"`
	DNI       string    `gorm:"type:varchar(20);unique;not null" json:"dni"`
	CreatedAt time.Time `json:"created_at"`
	Tickets   []Ticket  `gorm:"foreignKey:UserID" json:"tickets,omitempty"`
}

type Venue struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre          string    `gorm:"type:varchar(150);not null" json:"nombre"`
	Direccion       string    `gorm:"type:varchar(255);not null" json:"direccion"`
	Filas           int       `gorm:"not null" json:"filas"`
	ColumnasPorFila int       `gorm:"not null" json:"columnas_por_fila"`
	Capacidad       int       `gorm:"not null" json:"capacidad"`
	CreatedAt       time.Time `json:"created_at"`
}

type Event struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo      string    `gorm:"type:varchar(150);not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	Categoria   string    `gorm:"type:varchar(100);not null;index" json:"categoria"`
	Fecha       time.Time `gorm:"not null" json:"fecha"`
	Lugar       string    `gorm:"type:varchar(150);not null" json:"lugar"`
	Precio      float64   `gorm:"type:decimal(10,2);not null" json:"precio"`
	Imagen      string    `gorm:"type:varchar(500)" json:"imagen"`
	CupoMaximo  int       `gorm:"not null" json:"cupo_maximo"`
	CupoDispon  int       `gorm:"column:cupo_disponible;not null" json:"cupo_disponible"`
	VenueID     uint      `gorm:"not null;default:0" json:"venue_id"`
	Venue       *Venue    `gorm:"foreignKey:VenueID" json:"venue,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type Seat struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	EventID uint   `gorm:"not null;index" json:"event_id"`
	Fila    string `gorm:"type:varchar(5);not null" json:"fila"`
	Numero  int    `gorm:"not null" json:"numero"`
	Ocupado bool   `gorm:"default:false" json:"ocupado"`
	Event   *Event `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;" json:"-"`
}

type Ticket struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	EventID   uint      `gorm:"not null" json:"event_id"`
	SeatID    *uint     `json:"seat_id,omitempty"`
	FechaComp time.Time `json:"fecha_compra"`
	Estado    string    `gorm:"type:varchar(50);default:'activo'" json:"estado"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Event     *Event    `gorm:"foreignKey:EventID;constraint:OnDelete:RESTRICT;" json:"event,omitempty"`
	Seat      *Seat     `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
}