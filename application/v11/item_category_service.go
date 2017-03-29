package v11

import (
	"strconv"
	"strings"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// ItemCategoryService used to handle application logic related to Item Category resource.
type ItemCategoryService struct {
	ItemCategoryRepository  ItemCategoryRepositoryInterface
	ItemCategoryTransformer ItemCategoryTransformerInterface
	GrocerRepository        GrocerRepositoryInterface
	GrocerService           GrocerServiceInterface
	DealRepository          DealRepositoryInterface
	DealService             DealServiceInterface
}

// GetItemCategoryNames function used to retrieve name for all item category in database.
func (ics *ItemCategoryService) GetItemCategoryNames() ([]string, int) {
	itemCategories, totalItemCategory := ics.ItemCategoryRepository.GetAllCategoryNames()

	return itemCategories, totalItemCategory
}

// GetItemCategories function used to retrieve all item categories in database.
func (ics *ItemCategoryService) GetItemCategories() ([]*ItemCategory, int) {
	itemCategories, totalItemCategory := ics.ItemCategoryRepository.GetAll()

	return itemCategories, totalItemCategory
}

// GetItemCategoryByGUID function used to retrieve item category by GUID.
func (ics *ItemCategoryService) GetItemCategoryByGUID(guid string) *ItemCategory {
	return ics.ItemCategoryRepository.GetByGUID(guid)
}

// GetItemCategoryByID function used to retrieve item category by ID.
func (ics *ItemCategoryService) GetItemCategoryByID(itemCategoryID int) *ItemCategory {
	return ics.ItemCategoryRepository.GetByID(itemCategoryID)
}

func (ics *ItemCategoryService) TransformItemCategories(data interface{}, totalData int) *ItemCategoryResponse {
	return ics.ItemCategoryTransformer.TransformCollection(data, totalData)
}

func (ics *ItemCategoryService) GetGrocerCategoriesThoseHaveDealsIncludingDeals(userGUID, grocerGUID,
	latitude, longitude, dealLimitPerCategory, relations string) ([]*ItemCategory, *systems.ErrorData) {

	_, error := ics.GrocerService.CheckGrocerExistOrNotByGUID(grocerGUID)

	if error != nil {
		return nil, error
	}

	_, error = ics.GrocerService.CheckGrocerPublishOrNotByGUID(grocerGUID)

	if error != nil {
		return nil, error
	}

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	grocer := ics.GrocerRepository.GetByGUID(grocerGUID, "")

	categories := ics.ItemCategoryRepository.GetGrocerCategoriesForThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, grocer.ID,
		currentDateInGMT8, latitudeInFLoat64, longitudeInFLoat64)

	for key, category := range categories {
		deals, _ := ics.DealRepository.GetDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuotaAndCategory(userGUID, category.GUID,
			grocer.ID, latitudeInFLoat64, longitudeInFLoat64, currentDateInGMT8, "1", dealLimitPerCategory, "")

		deals = ics.DealService.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

		categories[key].Deals = deals
	}

	return categories, nil
}
