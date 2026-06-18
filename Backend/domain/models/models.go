package domain

import (
	"time"
)

// User representa a los usuarios del sistema (tanto Clientes como Administradores)
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(191);unique;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Nombre    string    `gorm:"type:varchar(100);not null" json:"nombre"`
	Apellido  string    `gorm:"type:varchar(100);not null" json:"apellido"`
	Rol       string    `gorm:"type:varchar(50);not null;default:'cliente'" json:"rol"` 
	CreatedAt time.Time `json:"created_at"`
	Tickets   []Ticket  `gorm:"foreignKey:UserID" json:"tickets,omitempty"`
}

// Event representa los eventos disponibles en el catálogo (Actualizado con Categoría)
type Event struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Titulo      string    `gorm:"type:varchar(150);not null" json:"titulo"`
	Descripcion string    `gorm:"type:text" json:"descripcion"`
	
	// NUEVO CAMPO: Para filtrar por tipo de evento (ej: 'Recitales', 'Teatro', 'Deportes')
	Categoria   string    `gorm:"type:varchar(100);not null;index" json:"categoria"` 
	
	Fecha       time.Time `gorm:"not null" json:"fecha"`
	Lugar       string    `gorm:"type:varchar(150);not null" json:"lugar"`
	Precio      float64   `gorm:"type:decimal(10,2);not null" json:"precio"`
	CupoMaximo  int       `gorm:"not null" json:"cupo_maximo"`
	CupoDispon  int       `gorm:"not null" json:"cupo_disponible"`
	CreatedAt   time.Time `json:"created_at"`
}

// Ticket representa las entradas compradas por los usuarios para un evento
type Ticket struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	EventID   uint      `gorm:"not null" json:"event_id"`
	FechaComp time.Time `json:"fecha_compra"`
	Estado    string    `gorm:"type:varchar(50);default:'activo'" json:"estado"`
	User      *User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	Event     *Event    `gorm:"foreignKey:EventID;constraint:OnDelete:RESTRICT;" json:"event,omitempty"`
}