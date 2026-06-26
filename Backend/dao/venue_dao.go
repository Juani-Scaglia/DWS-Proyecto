package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func GetAllVenues() ([]models.Venue, error) {
	if DB == nil {
		return nil, errors.New(errDBNula)
	}
	var venues []models.Venue
	err := DB.Find(&venues).Error
	return venues, err
}

func GetVenueByID(id uint) (models.Venue, error) {
	if DB == nil {
		return models.Venue{}, errors.New(errDBNula)
	}
	var venue models.Venue
	err := DB.First(&venue, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return venue, errors.New("establecimiento no encontrado")
	}
	return venue, err
}

func CreateVenue(venue *models.Venue) error {
	if DB == nil {
		return errors.New(errDBNula)
	}
	return DB.Create(venue).Error
}

func UpdateVenue(id uint, fields map[string]interface{}) error {
	if DB == nil {
		return errors.New(errDBNula)
	}
	return DB.Model(&models.Venue{}).Where("id = ?", id).Updates(fields).Error
}

func DeleteVenue(id uint) error {
	if DB == nil {
		return errors.New(errDBNula)
	}
	return DB.Delete(&models.Venue{}, id).Error
}
