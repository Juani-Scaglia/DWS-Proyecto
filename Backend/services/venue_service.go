package services

import (
	"errors"

	"backend/dao"
	domain "backend/domain/models"
)

type VenueInput struct {
	Nombre          string `json:"nombre" binding:"required"`
	Direccion       string `json:"direccion" binding:"required"`
	Filas           int    `json:"filas" binding:"required,gt=0,max=26"`
	ColumnasPorFila int    `json:"columnas_por_fila" binding:"required,gt=0"`
}

func GetAllVenues() ([]domain.Venue, error) {
	return dao.GetAllVenues()
}

func GetVenueByID(id uint) (domain.Venue, error) {
	return dao.GetVenueByID(id)
}

func CreateVenue(input VenueInput) (*domain.Venue, error) {
	venue := &domain.Venue{
		Nombre:          input.Nombre,
		Direccion:       input.Direccion,
		Filas:           input.Filas,
		ColumnasPorFila: input.ColumnasPorFila,
		Capacidad:       input.Filas * input.ColumnasPorFila,
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

	fields := map[string]interface{}{
		"nombre":            input.Nombre,
		"direccion":         input.Direccion,
		"filas":             input.Filas,
		"columnas_por_fila": input.ColumnasPorFila,
		"capacidad":         input.Filas * input.ColumnasPorFila,
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
