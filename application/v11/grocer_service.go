package v11

import (
	"strconv"
	"strings"
	"time"

	"fmt"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// GrocerService will handle all application logic related to Grocer resource.
type GrocerService struct {
	GrocerRepository GrocerRepositoryInterface
	DealRepository   DealRepositoryInterface
}

// CheckGrocerPublishOrNotByGUID function used to check grocer publish or not by Grocer GUID.
func (gs *GrocerService) CheckGrocerPublishOrNotByGUID(grocerGUID string) (*Grocer, *systems.ErrorData) {
	grocer := gs.GrocerRepository.GetByGUIDAndStatus(grocerGUID, "publish", "")

	if grocer.GUID == "" {
		return nil, Error.GenericError("403", systems.GrocerNotPublish, fmt.Sprintf(systems.TitleGrocerNotPublish, grocerGUID),
			"message", systems.ErrorGrocerNotPublish)
	}

	return grocer, nil
}

// CheckGrocerExistOrNotByGUID function used to check grocer exist or not by checking
// grocer GUID.
func (gs *GrocerService) CheckGrocerExistOrNotByGUID(grocerGUID string) (*Grocer, *systems.ErrorData) {
	grocer := gs.GrocerRepository.GetByGUID(grocerGUID, "")

	if grocer.GUID == "" {
		return nil, Error.ResourceNotFoundError("Grocer", "guid", grocerGUID)
	}

	return grocer, nil
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
