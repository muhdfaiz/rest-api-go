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
	RemoveDealCashbackAndSetItemDealExpired(userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData
	ViewDealDetails(dealGUID string, relations string) *Ads
	GetAvailableDealsForGuestUser(latitude string, longitude string, pageNumber string, pageLimit string, relations string) ([]*Deal, int)
	GetAvailableDealsForRegisteredUser(userGUID string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
	GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID string, latitude string, longitude string, dealLimitPerCategory string, relations string) []*ItemCategory
	GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID string, categoryGUID string, latitude string, longitude string,
		dealLimitPerSubcategory string, relations string) []*ItemSubCategory
	GetAvailableDealsByCategoryForRegisteredUser(userGUID string, category string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
	GetAvailableDealsForSubCategoryForRegisteredUser(userGUID string, category string, latitude string, longitude string, pageNumber string, pageLimit string,
		relations string) ([]*Deal, int)
}

type DealService struct {
	DealRepository            DealRepositoryInterface
	LocationService           location.LocationServiceInterface
	DealCashbackFactory       DealCashbackFactoryInterface
	ShoppingListItemFactory   ShoppingListItemFactoryInterface
	DealCashbackRepository    DealCashbackRepositoryInterface
	ItemRepository            ItemRepositoryInterface
	ItemCategoryService       ItemCategoryServiceInterface
	ItemSubCategoryRepository ItemSubCategoryRepositoryInterface
	GrocerRepository          GrocerRepositoryInterface
	DealCashbackService       DealCashbackServiceInterface
}

func (ds *DealService) GetDealsBasedOnUserShoppingListItem(userGUID string, shoppingListItem *ShoppingListItem,
	latitude string, longitude string, dealsCollection []*Deal) []*Deal {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")
	currentTimeInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("15:04")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, _ := ds.DealRepository.GetAllDealsForCategoryWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, shoppingListItem.Category,
		latitude1InFLoat64, longitude1InFLoat64, currentDateInGMT8, "1", "10000", "")

	filteredDealsUniqueForEachShoppingList := []*Deal{}

	if len(deals) < 1 {
		return nil
	}

	for _, deal := range deals {
		dealAlreadyExistInOtherItem := false

		for _, dealCollection := range dealsCollection {
			if deal.GUID == dealCollection.GUID {
				dealAlreadyExistInOtherItem = true
				break
			}
		}

		if dealAlreadyExistInOtherItem == false {
			filteredDealsUniqueForEachShoppingList = append(filteredDealsUniqueForEachShoppingList, deal)
		}
	}

	if len(filteredDealsUniqueForEachShoppingList) < 1 {
		return nil
	}

	filteredDealsByStartAndEndTime := []*Deal{}

	// Filtered deal those only has valid time
	for _, filteredDealUniqueForEachShoppingList := range filteredDealsUniqueForEachShoppingList {
		if filteredDealUniqueForEachShoppingList.Time != "" {
			// Example: 08:00-10:00;15:00-18:00
			dealTimeRanges := strings.Split(filteredDealUniqueForEachShoppingList.Time, ";")

			for _, dealTimeRange := range dealTimeRanges {
				// Example: 08:00-10:00
				dealTimes := strings.Split(dealTimeRange, "-")

				// If deals time was valid time add deal to deals slice
				if currentTimeInGMT8 >= dealTimes[0] && currentTimeInGMT8 < dealTimes[1] {
					filteredDealsByStartAndEndTime = append(filteredDealsByStartAndEndTime, filteredDealUniqueForEachShoppingList)
					break
				}
			}
		} else {
			filteredDealsByStartAndEndTime = append(filteredDealsByStartAndEndTime, filteredDealUniqueForEachShoppingList)
		}
	}

	if len(filteredDealsByStartAndEndTime) < 1 {
		return nil
	}

	filteredDealsByPositiveTags := []*Deal{}

	itemNameInLowercase := strings.ToLower(shoppingListItem.Name)
	splitItemNames := strings.Fields(itemNameInLowercase)

	// If positive_tag not empty, filtered deal those only has match positive_tag
	for _, filteredDealByStartAndEndTime := range filteredDealsByStartAndEndTime {
		if filteredDealByStartAndEndTime.PositiveTag != "" {
			for key := range splitItemNames {
				matchPositiveTag := strings.Contains(strings.ToLower(filteredDealByStartAndEndTime.PositiveTag), splitItemNames[key])

				if matchPositiveTag == true {
					filteredDealsByPositiveTags = append(filteredDealsByPositiveTags, filteredDealByStartAndEndTime)
					break
				}
			}
		} else {
			filteredDealsByPositiveTags = append(filteredDealsByPositiveTags, filteredDealByStartAndEndTime)
		}
	}

	if len(filteredDealsByPositiveTags) < 1 {
		return nil
	}

	filteredDealsByNegativeTags := []*Deal{}

	// If negative_tag not empty, filtered deal those only has match negative_tag
	for _, filteredDealByPositiveTags := range filteredDealsByPositiveTags {
		if filteredDealByPositiveTags.NegativeTag != "" {

			for key := range splitItemNames {
				matchNegativeTag := strings.Contains(strings.ToLower(filteredDealByPositiveTags.NegativeTag), splitItemNames[key])

				if matchNegativeTag == true {
					break
				}

				if len(filteredDealsByNegativeTags) < 3 {
					filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, filteredDealByPositiveTags)
				}
			}
		} else {
			if len(filteredDealsByNegativeTags) < 3 {
				filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, filteredDealByPositiveTags)
			}
		}
	}

	if len(filteredDealsByNegativeTags) < 1 {
		return nil
	}

	return filteredDealsByNegativeTags
}

func (ds *DealService) RemoveDealCashbackAndSetItemDealExpired(userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := ds.DealRepository.GetDealByGUIDAndValidStartEndDate(dealGUID, currentDateInGMT8)

	// If deal already expired
	if deal.GUID == "" {
		// Delete user deal cashback
		err := ds.DealCashbackFactory.DeleteByUserGUIDShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID)

		if err != nil {
			return err
		}

		// Shopping list item deal expired data
		data := map[string]interface{}{"deal_expired": 1}

		// Set user shopping list item deal expired to 1(true)
		err = ds.ShoppingListItemFactory.UpdateByUserGUIDShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID, data)

		if err != nil {
			return err
		}
	}

	return nil
}

// ViewDealDetails function used to retrieve deal details including the relations
func (ds *DealService) ViewDealDetails(dealGUID string, relations string) *Ads {
	// Retrieve deal ID
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
func (ds *DealService) GetAvailableDealsForRegisteredUser(userGUID string, latitude string, longitude string, pageNumber string, pageLimit string, relations string) ([]*Deal, int) {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitudeInFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals, totalDeal := ds.DealRepository.GetAllDealsWithinValidRangeStartDateEndDateUserLimitAndQuota(userGUID, latitude1InFLoat64, longitudeInFLoat64,
		currentDateInGMT8, pageNumber, pageLimit, relations)

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
		}
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
		}
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
	}

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
	}

	return deals, totalDeal
}
