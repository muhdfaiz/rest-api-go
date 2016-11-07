package v1

import "github.com/jinzhu/gorm"

type GrocerLocationRepositoryInterface interface {
	GetByIDLatitudeAndLongitude(id int, latitude float64, longitude float64) *GrocerLocation
}

type GrocerLocationRepository struct {
	DB *gorm.DB
}

func (glr *GrocerLocationRepository) GetByIDLatitudeAndLongitude(id int, latitude float64, longitude float64) *GrocerLocation {
	grocerLocation := &GrocerLocation{}

	glr.DB.Model(&GrocerLocation{}).Where(&GrocerLocation{ID: id, Lat: latitude, Lng: longitude}).Find(&grocerLocation)

	return grocerLocation
}
