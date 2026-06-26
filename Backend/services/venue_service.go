package services

import (
	"errors"
	"math"

	"backend/dao"
	domain "backend/domain/models"
)

type VenueInput struct {
	Nombre         string `json:"nombre" binding:"required"`
	Direccion      string `json:"direccion" binding:"required"`
	Tipo           string `json:"tipo"`
	Capacidad      int    `json:"capacidad"`
	CapTribunaNorte int   `json:"cap_tribuna_norte"`
	CapTribunaSur   int   `json:"cap_tribuna_sur"`
	CapTribunaEste  int   `json:"cap_tribuna_este"`
	CapTribunaOeste int   `json:"cap_tribuna_oeste"`
	CapCampo        int   `json:"cap_campo"`
}

func (v VenueInput) totalCapacidad() int {
	if v.Tipo == "escenario" {
		return v.Capacidad
	}
	return v.CapTribunaNorte + v.CapTribunaSur + v.CapTribunaEste + v.CapTribunaOeste + v.CapCampo
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

func normalizeTipo(t string) string {
	if t == "escenario" {
		return "escenario"
	}
	return "estadio"
}

func CreateVenue(input VenueInput) (*domain.Venue, error) {
	total := input.totalCapacidad()
	if total <= 0 {
		return nil, errors.New("la capacidad total debe ser mayor a 0")
	}
	filas, cols := calcularGrilla(total)
	tipo := normalizeTipo(input.Tipo)
	venue := &domain.Venue{
		Nombre:          input.Nombre,
		Direccion:       input.Direccion,
		Tipo:            tipo,
		Capacidad:       total,
		Filas:           filas,
		ColumnasPorFila: cols,
		CapPlateaNorte:  input.CapTribunaNorte,
		CapPlateaSur:    input.CapTribunaSur,
		CapTribunaEste:  input.CapTribunaEste,
		CapTribunaOeste: input.CapTribunaOeste,
		CapCampo:        input.CapCampo,
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
	total := input.totalCapacidad()
	if total <= 0 {
		return nil, errors.New("la capacidad total debe ser mayor a 0")
	}
	filas, cols := calcularGrilla(total)
	tipoUpd := normalizeTipo(input.Tipo)
	fields := map[string]interface{}{
		"nombre":                 input.Nombre,
		"direccion":              input.Direccion,
		"tipo":                   tipoUpd,
		"capacidad":              total,
		"filas":                  filas,
		"columnas_por_fila":      cols,
		"cap_platea_norte":       input.CapTribunaNorte,
		"cap_platea_sur":         input.CapTribunaSur,
		"cap_tribuna_este":       input.CapTribunaEste,
		"cap_tribuna_oeste":      input.CapTribunaOeste,
		"cap_platea_preferencial": 0,
		"cap_campo":              input.CapCampo,
	}
	if err := dao.UpdateVenue(id, fields); err != nil {
		return nil, err
	}
	venue, err := dao.GetVenueByID(id)
	if err != nil {
		return nil, err
	}

	var events []domain.Event
	dao.DB.Where("venue_id = ?", id).Find(&events)
	sectores := VenueSectores(venue)
	for _, ev := range events {
		tx := dao.DB.Begin()
		tx.Where("event_id = ?", ev.ID).Delete(&domain.Seat{})
		if err := dao.CreateSeatsForEvent(tx, ev.ID, sectores); err != nil {
			tx.Rollback()
			continue
		}
		tx.Model(&domain.Event{}).Where("id = ?", ev.ID).Updates(map[string]interface{}{
			"cupo_maximo":     total,
			"cupo_disponible": total,
		})
		tx.Commit()
	}

	return &venue, nil
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
	if v.Tipo == "escenario" {
		return []dao.SectorDef{{Nombre: "General", Capacidad: v.Capacidad}}
	}
	var sectores []dao.SectorDef
	if v.CapPlateaNorte > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Norte", Capacidad: v.CapPlateaNorte})
	}
	if v.CapPlateaSur > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Sur", Capacidad: v.CapPlateaSur})
	}
	if v.CapTribunaEste > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Este", Capacidad: v.CapTribunaEste})
	}
	if v.CapTribunaOeste > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Tribuna Oeste", Capacidad: v.CapTribunaOeste})
	}
	if v.CapCampo > 0 {
		sectores = append(sectores, dao.SectorDef{Nombre: "Campo", Capacidad: v.CapCampo})
	}
	return sectores
}
