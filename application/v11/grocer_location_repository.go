package v11

import "github.com/jinzhu/gorm"

// GrocerLocationRepository will handle all CRUD function for Grocer Location resource.
type GrocerLocationRepository struct {
	DB *gorm.DB
}

// GetByIDLatitudeAndLongitude function used to retrieve multiple grocer locatio.
// Filter by:
// - ID
// - latitude
// - longitude
func (glr *GrocerLocationRepository) GetByIDLatitudeAndLongitude(id int, latitude float64, longitude float64) *GrocerLocation {
	grocerLocation := &GrocerLocation{}

	glr.DB.Model(&GrocerLocation{}).Where(&GrocerLocation{ID: id, Lat: latitude, Lng: longitude}).Find(&grocerLocation)

	return grocerLocation
}
