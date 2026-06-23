package services

import (
	"errors"
	"math"

	"backend/dao"
	domain "backend/domain/models"
)

type VenueInput struct {
	Nombre    string `json:"nombre" binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Capacidad int    `json:"capacidad" binding:"required,gt=0"`
}

func calcularGrilla(capacidad int) (int, int) {
	cols := int(math.Ceil(math.Sqrt(float64(capacidad))))
	if cols > 50 {
		cols = 50
	}
	if cols < 1 {
		cols = 1
	}
	filas := int(math.Ceil(float64(capacidad) / float64(cols)))
	return filas, cols
}

func GetAllVenues() ([]domain.Venue, error) {
	return dao.GetAllVenues()
}

func GetVenueByID(id uint) (domain.Venue, error) {
	return dao.GetVenueByID(id)
}

func CreateVenue(input VenueInput) (*domain.Venue, error) {
	filas, cols := calcularGrilla(input.Capacidad)
	venue := &domain.Venue{
		Nombre:          input.Nombre,
		Direccion:       input.Direccion,
		Filas:           filas,
		ColumnasPorFila: cols,
		Capacidad:       input.Capacidad,
	}
	if err := dao.CreateVenue(venue); err != nil {
		return nil, err
	}
	return venue, nil
}

func UpdateVenue(id uint, input VenueInput) (*domain.Venue, error) {
	if _, err := dao.GetVenueByID(id); err != nil {
		return nil, err
	}

	var count int64
	dao.DB.Model(&domain.Event{}).Where("venue_id = ?", id).Count(&count)
	if count > 0 {
		return nil, errors.New("no se puede modificar un establecimiento que tiene eventos asociados")
	}

	filas, cols := calcularGrilla(input.Capacidad)
	fields := map[string]interface{}{
		"nombre":            input.Nombre,
		"direccion":         input.Direccion,
		"filas":             filas,
		"columnas_por_fila": cols,
		"capacidad":         input.Capacidad,
	}
	if err := dao.UpdateVenue(id, fields); err != nil {
		return nil, err
	}
	venue, err := dao.GetVenueByID(id)
	return &venue, err
}

func DeleteVenue(id uint) error {
	if _, err := dao.GetVenueByID(id); err != nil {
		return err
	}

	var count int64
	dao.DB.Model(&domain.Event{}).Where("venue_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("no se puede eliminar un establecimiento que tiene eventos asociados")
	}

	return dao.DeleteVenue(id)
}
