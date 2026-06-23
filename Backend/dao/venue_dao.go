package dao

import (
	"backend/domain/models"
	"errors"

	"gorm.io/gorm"
)

func GetAllVenues() ([]models.Venue, error) {
	var venues []models.Venue
	return venues, DB.Find(&venues).Error
}

func GetVenueByID(id uint) (models.Venue, error) {
	var venue models.Venue
	err := DB.First(&venue, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return venue, errors.New("establecimiento no encontrado")
	}
	return venue, err
}

func CreateVenue(venue *models.Venue) error {
	return DB.Create(venue).Error
}

func UpdateVenue(id uint, fields map[string]interface{}) error {
	return DB.Model(&models.Venue{}).Where("id = ?", id).Updates(fields).Error
}

func DeleteVenue(id uint) error {
	return DB.Delete(&models.Venue{}, id).Error
}
