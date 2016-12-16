package v1

import (
	"strconv"
	"strings"
	"time"
)

// GrocerServiceInterface is a contract that defines the method needed for Grocer Service.
type GrocerServiceInterface interface {
	GetGrocerByID(grocerID int, relations string) *Grocer
	GetAllGrocers(pageNumber, pageLimit, relations string) ([]*Grocer, int)
	GetAllGrocersIncludingDeals(userGUID, latitude, longitude string) []*Grocer
}

type GrocerService struct {
	GrocerRepository GrocerRepositoryInterface
	DealRepository   DealRepositoryInterface
}

// GetGrocerByID function used to retrieve grocer by grocer ID.
func (gs *GrocerService) GetGrocerByID(grocerID int, relations string) *Grocer {
	grocer := gs.GrocerRepository.GetByID(grocerID, relations)

	return grocer
}

// GetAllGrocers function used to retrieve all grocers and total number of grocers.
func (gs *GrocerService) GetAllGrocers(pageNumber, pageLimit, relations string) ([]*Grocer, int) {
	grocers, totalGrocer := gs.GrocerRepository.GetAll(pageNumber, pageLimit, relations)

	return grocers, totalGrocer
}

// GetAllGrocersIncludingDeals function used retrieve only grocer those have deals incuding
// deals related to the grocers.
func (gs *GrocerService) GetAllGrocersIncludingDeals(userGUID, latitude, longitude string) []*Grocer {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	grocers := gs.GrocerRepository.GetAllGrocersThoseOnlyHaveDeal()

	for key, grocer := range grocers {
		totalDeals := gs.DealRepository.CountDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, grocer.ID, latitudeInFLoat64, longitudeInFLoat64,
			currentDateInGMT8)

		grocers[key].TotalDeals = totalDeals
	}

	return grocers
}
