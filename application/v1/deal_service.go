package v1

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealService struct {
	DealRepository             DealRepositoryInterface
	DealTransformer            DealTransformerInterface
	LocationService            location.LocationServiceInterface
	DealCashbackRepository     DealCashbackRepositoryInterface
	ShoppingListItemRepository ShoppingListItemRepositoryInterface
	ItemRepository             ItemRepositoryInterface
	ItemCategoryRepository     ItemCategoryRepositoryInterface
	ItemSubCategoryRepository  ItemSubCategoryRepositoryInterface
	GrocerService              GrocerServiceInterface
	ShoppingListService        ShoppingListServiceInterface
}

// CheckDealExistOrNotByGUID function used to check deal exist or not by checking the deal GUID.
func (ds *DealService) CheckDealExistOrNotByGUID(dealGUID string) (*Deal, *systems.ErrorData) {
	deal := ds.DealRepository.GetDealByGUID(dealGUID)

	if deal.GUID == "" {
		return nil, Error.ResourceNotFoundError("Deal", "guid", dealGUID)
	}

	return deal, nil
}

// GetDealsBasedOnUserShoppingListItem function used to retrieve deals for each of user shopping list items.
// Maximum 3 deals for each of shopping list items.
// The deal is valid if item category must be same with shopping list item category.
// The deal is valid if user location within valid range (10KM radius).
// The deal is valid if today date is within deal start date and end date.
// The deal is valid if total number of deal added to list by user not exceed deal perlimit.
// The deal is valid when total deal added to list by all user below the deal quota.
// If deal start time and end time not empty, the deal is valid if current time within start time and end time.
// If deal positive tag not empty, the deal is valid if the item name contain any of the positive tag keyword.
// If deal negative tag not empty, the deal is not valid if the item name contain any of the negative tag keyword.
func (ds *DealService) GetDealsBasedOnUserShoppingListItem(userGUID, shoppingListGUID string, shoppingListItem *ShoppingListItem,
	latitude, longitude string) []*Deal {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")
	currentTimeInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("15:04")

	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, _ := ds.DealRepository.GetDealsBySubcategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, shoppingListItem.SubCategory,
		latitude1InFLoat64, longitude1InFLoat64, currentDateInGMT8, "1", "", "Category")

	if len(deals) < 1 {
		return nil
	}

	filteredDealsByStartAndEndTime := ds.FilteredDealMustBeWithinStartAndEndTime(deals, currentDateInGMT8, currentTimeInGMT8)

	if len(filteredDealsByStartAndEndTime) < 1 {
		return nil
	}

	filteredDealsByPositiveTags := ds.FilteredDealByPositiveTag(filteredDealsByStartAndEndTime, shoppingListItem.Name)

	if len(filteredDealsByPositiveTags) < 1 {
		return nil
	}

	filteredDealsByNegativeTags := ds.FilteredDealByNegativeTag(filteredDealsByPositiveTags, shoppingListItem.Name)

	if len(filteredDealsByNegativeTags) < 1 {
		return nil
	}

	filteredDealsNotAddedToList := ds.FilteredDealsNotAddedTolist(filteredDealsByNegativeTags, userGUID, shoppingListGUID)

	if len(filteredDealsNotAddedToList) < 1 {
		return nil
	}

	firstThreeDeals := ds.GetFirstThreeDeals(filteredDealsNotAddedToList)

	return firstThreeDeals
}

// GetDealsBasedOnSampleShoppingListItem function used to retrieve deals for each of sample shopping list items.
// Maximum 3 deals for each of sample shopping list items.
// The deal is valid if item category must be same with shopping list item category.
// The deal is valid if user location within valid range (10KM radius).
// The deal is valid if today date is within deal start date and end date.
// The deal is valid when total deal added to list by all user below the deal quota.
// If deal start time and end time not empty, the deal is valid if current time within start time and end time.
// If deal positive tag not empty, the deal is valid if the item name contain any of the positive tag keyword.
// If deal negative tag not empty, the deal is not valid if the item name contain any of the negative tag keyword.
func (ds *DealService) GetDealsBasedOnSampleShoppingListItem(defaultShoppingListItem *DefaultShoppingListItem, latitude,
	longitude string) []*Deal {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")
	currentTimeInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("15:04")

	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals := []*Deal{}

	if latitude == "" || longitude == "" {
		deals, _ = ds.DealRepository.GetDealsForCategoryWithinDateRangeAndQuota(defaultShoppingListItem.Category, currentDateInGMT8, "1", "", "Category")
	} else {
		deals, _ = ds.DealRepository.GetDealsForCategoryWithinRangeAndDateRangeAndQuota(defaultShoppingListItem.Category,
			latitude1InFLoat64, longitude1InFLoat64, currentDateInGMT8, "1", "", "Category")
	}

	if len(deals) < 1 {
		return nil
	}

	filteredDealsByStartAndEndTime := ds.FilteredDealMustBeWithinStartAndEndTime(deals, currentDateInGMT8, currentTimeInGMT8)

	if len(filteredDealsByStartAndEndTime) < 1 {
		return nil
	}

	filteredDealsByPositiveTags := ds.FilteredDealByPositiveTag(filteredDealsByStartAndEndTime, defaultShoppingListItem.Name)

	if len(filteredDealsByPositiveTags) < 1 {
		return nil
	}

	filteredDealsByNegativeTags := ds.FilteredDealByNegativeTag(filteredDealsByPositiveTags, defaultShoppingListItem.Name)

	if len(filteredDealsByNegativeTags) < 1 {
		return nil
	}

	firstThreeDeals := ds.GetFirstThreeDeals(filteredDealsByNegativeTags)

	return firstThreeDeals
}

// FilteredDealMustBeUniquePerShoppingList function used to set the deal must be unique for each of shopping lists items.
func (ds *DealService) FilteredDealMustBeUniquePerShoppingList(deals []*Deal, dealsCollection []*Deal, userGUID string) []*Deal {
	filteredDealsUniqueForEachShoppingList := []*Deal{}

	for _, deal := range deals {
		dealAlreadyExistInOtherItem := false

		for _, dealCollection := range dealsCollection {
			if deal.GUID == dealCollection.GUID {
				dealAlreadyExistInOtherItem = true
				break
			}
		}

		if dealAlreadyExistInOtherItem == false {

			if userGUID != "" {
				deal = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeal(deal, userGUID)
			}

			filteredDealsUniqueForEachShoppingList = append(filteredDealsUniqueForEachShoppingList, deal)
		}
	}

	return filteredDealsUniqueForEachShoppingList
}

// FilteredDealMustBeWithinStartAndEndTime function used to find deal that still within the deal time.
func (ds *DealService) FilteredDealMustBeWithinStartAndEndTime(deals []*Deal, currentDateInGMT8, currentTimeInGMT8 string) []*Deal {
	filteredDealsByStartAndEndTime := []*Deal{}

	// Filtered deal those only has valid time
	for _, deal := range deals {
		if deal.Time != "" {
			// Example: 08:00-10:00;15:00-18:00
			dealTimeRanges := strings.Split(deal.Time, ";")

			for _, dealTimeRange := range dealTimeRanges {
				// Example: 08:00-10:00
				dealTimes := strings.Split(dealTimeRange, "-")

				// If deals time was valid time add deal to deals slice
				if currentTimeInGMT8 >= dealTimes[0] && currentTimeInGMT8 < dealTimes[1] {
					filteredDealsByStartAndEndTime = append(filteredDealsByStartAndEndTime, deal)
					break
				}
			}
		} else {
			filteredDealsByStartAndEndTime = append(filteredDealsByStartAndEndTime, deal)
		}
	}

	return filteredDealsByStartAndEndTime
}

// FilteredDealByPositiveTag function used to find deal that match positive tag with the deal item name.
func (ds *DealService) FilteredDealByPositiveTag(deals []*Deal, shoppingListItemName string) []*Deal {
	filteredDealsByPositiveTags := []*Deal{}

	itemNameInLowercase := strings.ToLower(shoppingListItemName)
	splitItemNames := strings.Fields(itemNameInLowercase)

	// If positive_tag not empty, filtered deal those only has match positive_tag
	for _, deal := range deals {
		if deal.PositiveTag != "" {
			for key := range splitItemNames {
				matchPositiveTag := strings.Contains(strings.ToLower(deal.PositiveTag), splitItemNames[key])

				if matchPositiveTag == true {
					filteredDealsByPositiveTags = append(filteredDealsByPositiveTags, deal)
					break
				}
			}
		} else {
			filteredDealsByPositiveTags = append(filteredDealsByPositiveTags, deal)
		}
	}

	return filteredDealsByPositiveTags
}

// FilteredDealByNegativeTag function used to find deal that don't match negative tag with the deal item name.
func (ds *DealService) FilteredDealByNegativeTag(deals []*Deal, shoppingListItemName string) []*Deal {
	filteredDealsByNegativeTags := []*Deal{}

	itemNameInLowercase := strings.ToLower(shoppingListItemName)
	splitItemNames := strings.Fields(itemNameInLowercase)

	// If negative_tag not empty, filtered deal those only has match negative_tag
	for _, deal := range deals {
		if deal.NegativeTag != "" {

			for key := range splitItemNames {
				matchNegativeTag := strings.Contains(strings.ToLower(deal.NegativeTag), splitItemNames[key])

				if matchNegativeTag == true {
					break
				}

				if len(filteredDealsByNegativeTags) < 3 {
					filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, deal)
				}
			}
		} else {
			if len(filteredDealsByNegativeTags) < 3 {
				filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, deal)
			}
		}
	}

	return filteredDealsByNegativeTags
}

// FilteredDealsNotAddedTolist will filter only deal not added to list by User for the shopping list.
func (ds *DealService) FilteredDealsNotAddedTolist(deals []*Deal, userGUID, shoppingListGUID string) []*Deal {
	filteredDealsNotAddedTolist := []*Deal{}

	for _, deal := range deals {

		dealCashback := ds.DealCashbackRepository.GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, deal.GUID)

		if dealCashback.GUID == "" {
			filteredDealsNotAddedTolist = append(filteredDealsNotAddedTolist, deal)
		}
	}

	return filteredDealsNotAddedTolist
}

// GetFirstThreeDeals will retrieve first 3 deal only from deal collection.
func (ds *DealService) GetFirstThreeDeals(deals []*Deal) []*Deal {
	firstThreeDeals := []*Deal{}

	for key, deal := range deals {
		if key <= 2 {
			firstThreeDeals = append(firstThreeDeals, deal)
		}
	}

	return firstThreeDeals
}

// RemoveDealCashbackAndSetItemDealExpired function used to soft delete deal cashback that already expired and set the item deal expired.
func (ds *DealService) RemoveDealCashbackAndSetItemDealExpired(dbTransaction *gorm.DB, userGUID, shoppingListGUID, dealGUID string) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := ds.DealRepository.GetDealByGUIDAndValidStartEndDate(dealGUID, currentDateInGMT8)

	if deal.GUID == "" {
		error := ds.DealCashbackRepository.DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(dbTransaction, userGUID, shoppingListGUID, dealGUID)

		if error != nil {
			return error
		}

		error = ds.ShoppingListItemRepository.SetDealExpired(dbTransaction, dealGUID)

		if error != nil {
			return error
		}
	}

	return nil
}

// ViewDealDetails function used to retrieve deal details including the relations
func (ds *DealService) ViewDealDetails(dealGUID, relations string) *Ads {
	deal := ds.DealRepository.GetDealByGUID(dealGUID)

	if deal.GUID == "" {
		return &Ads{}
	}

	dealWithRelations := ds.DealRepository.GetDealByIDWithRelations(deal.ID, relations)

	return dealWithRelations
}

// GetAvailableDealsForGuestUser function used to retrieve all deals within valid range 10KM
func (ds *DealService) GetAvailableDealsForGuestUser(latitude, longitude, pageNumber, pageLimit, relations string) ([]*Deal, int) {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetDealsWithinRangeAndDateRangeAndQuota(latitudeInFLoat64, longitudeInFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	for key, deal := range deals {
		deals[key].Items = ds.ItemRepository.GetByID(deal.ItemID, "Categories,Subcategories")
		deals[key].Grocerexclusives = ds.GrocerService.GetGrocerByID(deal.GrocerExclusive, "")
	}

	return deals, totalDeal
}

// GetAvailableDealsForRegisteredUser function used to retrieve all deals within valid range 10KM
func (ds *DealService) GetAvailableDealsForRegisteredUser(userGUID, name, latitude, longitude, pageNumber, pageLimit,
	relations string) ([]*Deal, int) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetDealsWithinRangeAndDateRangeAndUserLimitAndQuotaAndName(userGUID, name, latitude1InFLoat64, longitudeInFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetAvailableDealsGroupByCategoryForRegisteredUser function used to retrieve all deals group by category
func (ds *DealService) GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID, latitude, longitude, dealLimitPerCategory,
	relations string) []*ItemCategory {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	dealCategories, _ := ds.ItemCategoryRepository.GetAll()

	for key, dealCategory := range dealCategories {
		deals, totalDeal := ds.DealRepository.GetDealsByCategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, dealCategory.Name, latitude1InFLoat64, longitude1InFLoat64,
			currentDateInGMT8, "1", dealLimitPerCategory, relations)

		if len(deals) > 0 {
			dealCategories[key].Deals = deals
			dealCategories[key].TotalDeals = totalDeal
		}

		deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)
	}

	return dealCategories
}

// GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser function used to retrieve all deals group by category
func (ds *DealService) GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID, categoryGUID, latitude, longitude,
	dealLimitPerSubcategory, relations string) []*ItemSubCategory {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	uniqueSubCategories := ds.ItemSubCategoryRepository.GetSubCategoriesForCategoryThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, categoryGUID,
		currentDateInGMT8, latitudeInFLoat64, longitudeInFLoat64)

	for key, uniqueSubCategory := range uniqueSubCategories {
		deals, totalDeal := ds.DealRepository.GetDealBySubCategoryGUIDWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, uniqueSubCategory.GUID, latitudeInFLoat64, longitudeInFLoat64,
			currentDateInGMT8, "1", dealLimitPerSubcategory, relations)

		if len(deals) > 0 {
			uniqueSubCategories[key].Deals = deals
			uniqueSubCategories[key].TotalDeals = totalDeal
		}

		deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)
	}

	return uniqueSubCategories
}

// GetAvailableDealsByCategoryForRegisteredUser function used to retrieve all deals for specific category
func (ds *DealService) GetAvailableDealsByCategoryForRegisteredUser(userGUID, categoryName, latitude, longitude,
	pageNumber, pageLimit, relations string) ([]*Deal, int) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetDealsByCategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, categoryName, latitude1InFLoat64, longitude1InFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetAvailableDealsForSubCategoryForRegisteredUser function used to retrieve all deals for specific subcategory
func (ds *DealService) GetAvailableDealsForSubCategoryForRegisteredUser(userGUID, subCategoryGUID, latitude, longitude,
	pageNumber, pageLimit, relations string) ([]*Deal, int) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetDealBySubCategoryGUIDWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, subCategoryGUID,
		latitude1InFLoat64, longitude1InFLoat64, currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetAvailableDealsForGrocerByCategory function used to retrieve all deals for specific subcategory
func (ds *DealService) GetAvailableDealsForGrocerByCategory(request *http.Request, userGUID, grocerGUID, categoryGUID, latitude, longitude,
	pageNumber, pageLimit, relations string) (*DealResponse, *systems.ErrorData) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	grocer, error := ds.GrocerService.CheckGrocerExistOrNotByGUID(grocerGUID)

	if error != nil {
		return nil, error
	}

	itemCategory := ds.ItemCategoryRepository.GetByGUID(categoryGUID)

	if itemCategory.GUID == "" {
		return nil, error
	}

	deals, totalDeal := ds.DealRepository.GetDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuotaAndCategory(userGUID, categoryGUID,
		grocer.ID, latitudeInFLoat64, longitudeInFLoat64, currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	dealResponse := ds.DealTransformer.transformCollection(request, deals, totalDeal, pageLimit)

	return dealResponse, nil
}

// SetAddTolistInfoAndItemsAndGrocerExclusiveForDeal function used to set number of deal added to list by user,
// remaining time user can add to list, and can add to list status to tell user allow to add the deal or not.
func (ds *DealService) SetAddTolistInfoAndItemsAndGrocerExclusiveForDeal(deal *Deal, userGUID string) *Deal {
	if userGUID == "" {
		return deal
	}

	deal.CanAddTolist = 1

	totalNumberOfDealAddedToList := ds.DealCashbackRepository.CountByDealGUIDAndUserGUID(deal.GUID, userGUID)

	if totalNumberOfDealAddedToList >= deal.Perlimit {
		deal.CanAddTolist = 0
	}

	deal.NumberOfDealAddedToList = totalNumberOfDealAddedToList
	deal.RemainingAddToList = deal.Perlimit - totalNumberOfDealAddedToList
	deal.Items = ds.ItemRepository.GetByID(deal.ItemID, "Categories,Subcategories")
	deal.Grocerexclusives = ds.GrocerService.GetGrocerByID(deal.GrocerExclusive, "")

	return deal
}

// SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals function used to set number of deal added to list by user,
// remaining time user can add to list, and can add to list status to tell user allow to add the deal or not.
func (ds *DealService) SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals []*Deal, userGUID string) []*Deal {
	for key, deal := range deals {
		deals[key].CanAddTolist = 1

		totalNumberOfDealAddedToList := ds.DealCashbackRepository.CountByDealGUIDAndUserGUID(deal.GUID, userGUID)

		if totalNumberOfDealAddedToList >= deals[key].Perlimit {
			deals[key].CanAddTolist = 0
		}

		deals[key].NumberOfDealAddedToList = totalNumberOfDealAddedToList
		deals[key].RemainingAddToList = deal.Perlimit - totalNumberOfDealAddedToList
		deals[key].Items = ds.ItemRepository.GetByID(deal.ItemID, "Categories,Subcategories")
		deals[key].Grocerexclusives = ds.GrocerService.GetGrocerByID(deal.GrocerExclusive, "")
	}

	return deals
}
