package v1

import (
	"strconv"
	"strings"
)

type GrocerLocationServiceInterface interface {
	GetGrocersFromConvertionLocation(convertionLocation string) []*GrocerLocation
}

type GrocerLocationService struct {
	GrocerLocationRepository GrocerLocationRepositoryInterface
}

func (gls *GrocerLocationService) GetGrocersFromConvertionLocation(convertionLocation string) []*GrocerLocation {
	// Example: 4.1762919672722454,101.65968596935272;3.1661901679226436,101.69192880392075
	splitConversionLocations := strings.Split(convertionLocation, ";")

	grocerLocations := []*GrocerLocation{}

	for _, splitConversionLocation := range splitConversionLocations {
		// Example: 4.1762919672722454,101.65968596935272
		latitudeLongitudeAndGrocerID := strings.Split(splitConversionLocation, ",")

		// Convert Latitude and Longitude from string to float65
		latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitudeLongitudeAndGrocerID[0]), 64)
		longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitudeLongitudeAndGrocerID[1]), 64)

		// Convert grocer location ID from string to int
		grocerLocationID, _ := strconv.Atoi(latitudeLongitudeAndGrocerID[2])

		// Retrieve Grocer by latitude and longitude
		grocerLocation := gls.GrocerLocationRepository.GetByIDLatitudeAndLongitude(grocerLocationID, latitudeInFLoat64, longitudeInFLoat64)

		if grocerLocation.GUID != "" {
			grocerLocations = append(grocerLocations, grocerLocation)
		}
	}

	return grocerLocations
}
