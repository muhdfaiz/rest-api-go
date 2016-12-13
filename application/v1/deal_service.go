package v1

import (
	"strconv"
	"strings"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealServiceInterface interface {
	GetDealsBasedOnUserShoppingListItem(userGUID string, shopppingListItems *ShoppingListItem, latitude string,
		longitude string, dealsCollection []*Deal) []*Deal
	FilteredDealMustBeUniqueForEachOfShoppingListItem(deals []*Deal, dealsCollection []*Deal, userGUID string) []*Deal
	FilteredDealMustBeWithinStartAndEndTime(deals []*Deal, currentDateInGMT8 string, currentTimeInGMT8 string) []*Deal
	FilteredDealByPositiveTag(deals []*Deal, shoppingListItemName string) []*Deal
	RemoveDealCashbackAndSetItemDealExpired(userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData
	ViewDealDetails(dealGUID string, relations string) *Ads
	GetAvailableDealsForGuestUser(latitude string, longitude string, pageNumber string, pageLimit string, relations string) ([]*Deal, int)
	GetAvailableDealsForRegisteredUser(userGUID string, name string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
	GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID string, latitude string, longitude string, dealLimitPerCategory string, relations string) []*ItemCategory
	GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID string, categoryGUID string, latitude string, longitude string,
		dealLimitPerSubcategory string, relations string) []*ItemSubCategory
	GetAvailableDealsByCategoryForRegisteredUser(userGUID string, category string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
	GetAvailableDealsForSubCategoryForRegisteredUser(userGUID string, category string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
	GetDealByGUID(dealGUID string) *Deal
	SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals []*Deal, userGUID string) []*Deal
}

type DealService struct {
	DealRepository             DealRepositoryInterface
	LocationService            location.LocationServiceInterface
	DealCashbackFactory        DealCashbackFactoryInterface
	ShoppingListItemRepository ShoppingListItemRepositoryInterface
	DealCashbackRepository     DealCashbackRepositoryInterface
	ItemRepository             ItemRepositoryInterface
	ItemCategoryService        ItemCategoryServiceInterface
	ItemSubCategoryRepository  ItemSubCategoryRepositoryInterface
	GrocerRepository           GrocerRepositoryInterface
	DealCashbackService        DealCashbackServiceInterface
}

func (ds *DealService) GetDealsBasedOnUserShoppingListItem(userGUID string, shoppingListItem *ShoppingListItem,
	latitude string, longitude string, dealsCollection []*Deal) []*Deal {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")
	currentTimeInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("15:04")

	// Convert Latitude and Longitude from string to float64
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, _ := ds.DealRepository.GetAllDealsForCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, shoppingListItem.Category,
		latitude1InFLoat64, longitude1InFLoat64, currentDateInGMT8, "1", "10000", "Category")

	if len(deals) < 1 {
		return nil
	}

	filteredDealsUniqueForEachShoppingList := ds.FilteredDealMustBeUniqueForEachOfShoppingListItem(deals, dealsCollection, userGUID)

	if len(filteredDealsUniqueForEachShoppingList) < 1 {
		return nil
	}

	filteredDealsByStartAndEndTime := ds.FilteredDealMustBeWithinStartAndEndTime(filteredDealsUniqueForEachShoppingList, currentDateInGMT8, currentTimeInGMT8)

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

	return filteredDealsByNegativeTags
}

// FilteredDealMustBeUniqueForEachOfShoppingListItem function used to set the deal must be unique for each of shopping lists items.
func (ds *DealService) FilteredDealMustBeUniqueForEachOfShoppingListItem(deals []*Deal, dealsCollection []*Deal, userGUID string) []*Deal {
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
			deal.CanAddTolist = 1

			// Check If deal quota still available for the user.
			total := ds.DealCashbackRepository.CountByDealGUIDAndUserGUID(deal.GUID, userGUID)

			if total >= deal.Perlimit {
				deal.CanAddTolist = 0
			}

			deal.NumberOfDealAddedToList = total
			deal.RemainingAddToList = deal.Perlimit - total
			filteredDealsUniqueForEachShoppingList = append(filteredDealsUniqueForEachShoppingList, deal)
		}
	}

	return filteredDealsUniqueForEachShoppingList
}

// FilteredDealMustBeWithinStartAndEndTime function used to find deal that still within the deal time.
func (ds *DealService) FilteredDealMustBeWithinStartAndEndTime(deals []*Deal, currentDateInGMT8 string, currentTimeInGMT8 string) []*Deal {
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

// RemoveDealCashbackAndSetItemDealExpired function used to soft delete deal cashback that already expired and set the item deal expired.
func (ds *DealService) RemoveDealCashbackAndSetItemDealExpired(userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := ds.DealRepository.GetDealByGUIDAndValidStartEndDate(dealGUID, currentDateInGMT8)

	if deal.GUID == "" {
		error := ds.DealCashbackFactory.DeleteByUserGUIDShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID)

		if error != nil {
			return error
		}

		error = ds.ShoppingListItemRepository.SetDealExpired(dealGUID)

		if error != nil {
			return error
		}
	}

	return nil
}

// ViewDealDetails function used to retrieve deal details including the relations
func (ds *DealService) ViewDealDetails(dealGUID string, relations string) *Ads {
	deal := ds.DealRepository.GetDealByGUID(dealGUID)

	if deal.GUID == "" {
		return &Ads{}
	}

	dealWithRelations := ds.DealRepository.GetDealByIDWithRelations(deal.ID, relations)

	return dealWithRelations
}

// GetAvailableDealsForGuestUser function used to retrieve all deals within valid range 10KM
func (ds *DealService) GetAvailableDealsForGuestUser(latitude string, longitude string, pageNumber string, pageLimit string, relations string) ([]*Deal, int) {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetAllDealsWithinValidRangeStartDateEndDateAndQuota(latitudeInFLoat64, longitudeInFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	for key, deal := range deals {
		deals[key].Items = ds.ItemRepository.GetByID(deal.ItemID, "Categories,Subcategories")
		deals[key].Grocerexclusives = ds.GrocerRepository.GetByID(deal.GrocerExclusive, "")
	}

	return deals, totalDeal
}

// GetAvailableDealsForRegisteredUser function used to retrieve all deals within valid range 10KM
func (ds *DealService) GetAvailableDealsForRegisteredUser(userGUID string, name string, latitude string, longitude string, pageNumber string, pageLimit string, relations string) ([]*Deal, int) {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetAllDealsWithinValidRangeStartDateEndDateUserLimitQuotaAndName(userGUID, name, latitude1InFLoat64, longitudeInFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetAvailableDealsGroupByCategoryForRegisteredUser function used to retrieve all deals group by category
func (ds *DealService) GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID string, latitude string, longitude string, dealLimitPerCategory string, relations string) []*ItemCategory {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	uniqueDealCategories, _ := ds.ItemCategoryService.GetItemCategories()

	for key, uniqueDealCategory := range uniqueDealCategories {
		deals, totalDeal := ds.DealRepository.GetAllDealsForCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, uniqueDealCategory.Name, latitude1InFLoat64, longitude1InFLoat64,
			currentDateInGMT8, "1", dealLimitPerCategory, relations)

		if len(deals) > 0 {
			uniqueDealCategories[key].Deals = deals
			uniqueDealCategories[key].TotalDeals = totalDeal
		}

		deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)
	}

	return uniqueDealCategories
}

// GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser function used to retrieve all deals group by category
func (ds *DealService) GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID string, categoryGUID string, latitude string, longitude string,
	dealLimitPerSubcategory string, relations string) []*ItemSubCategory {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	uniqueSubCategories := ds.DealRepository.GetUniqueSubCategoriesForDealsWithinValidRangeStartDateEndDateUserLimitSubCategoryAndQuota(userGUID, categoryGUID,
		latitudeInFLoat64, longitudeInFLoat64, currentDateInGMT8, "", "", "")

	for key, uniqueSubCategory := range uniqueSubCategories {
		deals, totalDeal := ds.DealRepository.GetAllDealsForSubCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, uniqueSubCategory.GUID, latitudeInFLoat64, longitudeInFLoat64,
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
func (ds *DealService) GetAvailableDealsByCategoryForRegisteredUser(userGUID string, categoryName string, latitude string, longitude string,
	pageNumber string, pageLimit string, relations string) ([]*Deal, int) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetAllDealsForCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, categoryName, latitude1InFLoat64, longitude1InFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetAvailableDealsForSubCategoryForRegisteredUser function used to retrieve all deals for specific subcategory
func (ds *DealService) GetAvailableDealsForSubCategoryForRegisteredUser(userGUID string, subCategoryGUID string, latitude string, longitude string,
	pageNumber string, pageLimit string, relations string) ([]*Deal, int) {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetAllDealsForSubCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, subCategoryGUID, latitude1InFLoat64, longitude1InFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

	deals = ds.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals, userGUID)

	return deals, totalDeal
}

// GetDealByGUID function used to retrieve deal by deal GUID.
func (ds *DealService) GetDealByGUID(dealGUID string) *Deal {
	deal := ds.DealRepository.GetDealByGUID(dealGUID)

	return deal
}

// SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals function used to set number of deal added to list by user,
// remaining time user can add to list, and can add to list status to tell user allow to add the deal or not.
func (ds *DealService) SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals []*Deal, userGUID string) []*Deal {
	for key, deal := range deals {
		deals[key].CanAddTolist = 1

		// Check If deal quota still available for the user.
		total := ds.DealCashbackRepository.CountByDealGUIDAndUserGUID(deal.GUID, userGUID)

		if total >= deals[key].Perlimit {
			deals[key].CanAddTolist = 0
		}

		deals[key].NumberOfDealAddedToList = total
		deals[key].RemainingAddToList = deal.Perlimit - total
		deals[key].Items = ds.ItemRepository.GetByID(deal.ItemID, "Categories,Subcategories")
		deals[key].Grocerexclusives = ds.GrocerRepository.GetByID(deal.GrocerExclusive, "")
	}

	return deals
}
