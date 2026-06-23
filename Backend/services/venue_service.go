package services

import (
	"errors"
	"math"

	"backend/dao"
	domain "backend/domain/models"
)

type VenueInput struct {
	Nombre               string `json:"nombre" binding:"required"`
	Direccion            string `json:"direccion" binding:"required"`
	CapPlateaNorte       int    `json:"cap_platea_norte"`
	CapPlateaSur         int    `json:"cap_platea_sur"`
	CapTribunaEste       int    `json:"cap_tribuna_este"`
	CapTribunaOeste      int    `json:"cap_tribuna_oeste"`
	CapPlateaPreferencial int   `json:"cap_platea_preferencial"`
	CapCampo             int    `json:"cap_campo"`
}

func (v VenueInput) totalCapacidad() int {
	return v.CapPlateaNorte + v.CapPlateaSur + v.CapTribunaEste + v.CapTribunaOeste + v.CapPlateaPreferencial + v.CapCampo
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
	total := input.totalCapacidad()
	if total <= 0 {
		return nil, errors.New("la capacidad total debe ser mayor a 0")
	}
	filas, cols := calcularGrilla(total)
	venue := &domain.Venue{
		Nombre:                input.Nombre,
		Direccion:             input.Direccion,
		Capacidad:             total,
		Filas:                 filas,
		ColumnasPorFila:       cols,
		CapPlateaNorte:        input.CapPlateaNorte,
		CapPlateaSur:          input.CapPlateaSur,
		CapTribunaEste:        input.CapTribunaEste,
		CapTribunaOeste:       input.CapTribunaOeste,
		CapPlateaPreferencial: input.CapPlateaPreferencial,
		CapCampo:              input.CapCampo,
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
	total := input.totalCapacidad()
	if total <= 0 {
		return nil, errors.New("la capacidad total debe ser mayor a 0")
	}
	filas, cols := calcularGrilla(total)
	fields := map[string]interface{}{
		"nombre":                 input.Nombre,
		"direccion":              input.Direccion,
		"capacidad":              total,
		"filas":                  filas,
		"columnas_por_fila":      cols,
		"cap_platea_norte":       input.CapPlateaNorte,
		"cap_platea_sur":         input.CapPlateaSur,
		"cap_tribuna_este":       input.CapTribunaEste,
		"cap_tribuna_oeste":      input.CapTribunaOeste,
		"cap_platea_preferencial": input.CapPlateaPreferencial,
		"cap_campo":              input.CapCampo,
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

func VenueSectores(v domain.Venue) []dao.SectorDef {
	var sectores []dao.SectorDef
	if v.CapPlateaNorte > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Platea Norte", Capacidad: v.CapPlateaNorte})
	}
	if v.CapPlateaSur > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Platea Sur", Capacidad: v.CapPlateaSur})
	}
	if v.CapTribunaEste > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Este", Capacidad: v.CapTribunaEste})
	}
	if v.CapTribunaOeste > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Oeste", Capacidad: v.CapTribunaOeste})
	}
	if v.CapPlateaPreferencial > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Preferencial", Capacidad: v.CapPlateaPreferencial})
	}
	if v.CapCampo > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Campo", Capacidad: v.CapCampo})
	}
	return sectores
}
