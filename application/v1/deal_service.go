package v1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealServiceInterface interface {
	GetDealsBasedOnUserShoppingListItem(userGUID string, shopppingListItems *ShoppingListItem,
		latitude string, longitude string, dealsCollection []*Deal) []*Deal
	RemoveDealCashbackAndSetItemDealExpired(userGUID string, dealGUID string) *systems.ErrorData
	ViewDealDetails(dealGUID string, relations string) *Ads
	GetAllDealsBasedOnLatitudeAndLongitudeAndQuota(latitude string, longitude string, offset string, limit string, relations string) ([]*Deal, int)
}

type DealService struct {
	DealRepository          DealRepositoryInterface
	LocationService         location.LocationServiceInterface
	DealCashbackFactory     DealCashbackFactoryInterface
	ShoppingListItemFactory ShoppingListItemFactoryInterface
	DealCashbackRepository  DealCashbackRepositoryInterface
	ItemRepository          ItemRepositoryInterface
}

func (ds *DealService) GetDealsBasedOnUserShoppingListItem(userGUID string, shoppingListItem *ShoppingListItem,
	latitude string, longitude string, dealsCollection []*Deal) []*Deal {

	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")
	currentTimeInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("15:04")

	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	deals := ds.DealRepository.GetDealsByCategoryAndValidStartEndDate(currentDateInGMT8, shoppingListItem)

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

	// If positive_tag not empty, filtered deal those only has match positive_tag
	for _, filteredDealByStartAndEndTime := range filteredDealsByStartAndEndTime {
		if filteredDealByStartAndEndTime.PositiveTag != "" {
			splitItemNames := strings.Split(strings.ToLower(shoppingListItem.Name), " ")

			for _, splitItemName := range splitItemNames {
				matchPositiveTag := strings.Contains(strings.ToLower(filteredDealByStartAndEndTime.PositiveTag), splitItemName)

				if matchPositiveTag {
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
			splitItemNames := strings.Split(strings.ToLower(shoppingListItem.Name), " ")

			for _, splitItemName := range splitItemNames {
				matchNegativeTag := strings.Contains(strings.ToLower(filteredDealByPositiveTags.NegativeTag), splitItemName)

				if matchNegativeTag {
					break
				}

				filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, filteredDealByPositiveTags)
				break
			}
		} else {
			filteredDealsByNegativeTags = append(filteredDealsByNegativeTags, filteredDealByPositiveTags)
		}
	}

	if len(filteredDealsByNegativeTags) < 1 {
		return nil
	}

	filteredDealsByConversionLocation := []*Deal{}

	// If negative_tag not empty, filtered deal those only has match negative_tag
	for _, filteredDealByNegativeTags := range filteredDealsByNegativeTags {
		if filteredDealByNegativeTags.ConversionLocation != "" {
			splitConversionLocations := strings.Split(filteredDealByNegativeTags.ConversionLocation, ";")

			for _, splitConversionLocation := range splitConversionLocations {
				splitConversionLocations := strings.Split(splitConversionLocation, ",")

				// Convert Latitude and Longitude from string to float64
				latitude2InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(splitConversionLocations[0]), 64)
				longitude2InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(splitConversionLocations[1]), 64)

				distanceInKM := ds.LocationService.CalculateDistanceInKilometer(latitude1InFLoat64, longitude1InFLoat64, latitude2InFLoat64, longitude2InFLoat64)

				if distanceInKM < 10 {
					filteredDealsByConversionLocation = append(filteredDealsByConversionLocation, filteredDealByNegativeTags)
					break
				}
			}
		} else {
			filteredDealsByConversionLocation = append(filteredDealsByConversionLocation, filteredDealByNegativeTags)
		}
	}

	if len(filteredDealsByConversionLocation) < 1 {
		return nil
	}

	filteredDealsByQuota := []*Deal{}

	// Check deal already used more than quota available
	for _, filteredDealByConversionLocation := range filteredDealsByConversionLocation {
		totalDealAddedToList := ds.DealCashbackRepository.CountByDealGUID(filteredDealByConversionLocation.GUID)
		fmt.Println(totalDealAddedToList)
		if totalDealAddedToList < filteredDealByConversionLocation.Quota {
			filteredDealsByQuota = append(filteredDealsByQuota, filteredDealByConversionLocation)
		}
	}

	if len(filteredDealsByQuota) < 1 {
		return nil
	}

	return filteredDealsByQuota
}

func (ds *DealService) RemoveDealCashbackAndSetItemDealExpired(userGUID string, dealGUID string) *systems.ErrorData {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	deal := ds.DealRepository.GetDealByGUIDAndValidStartEndDate(dealGUID, currentDateInGMT8)

	// If deal already expired
	if deal.GUID == "" {
		// Delete user deal cashback
		err := ds.DealCashbackFactory.DeleteByUserGUIDAndDealGUID(userGUID, dealGUID)

		if err != nil {
			return err
		}

		// Shopping list item deal expired data
		data := map[string]interface{}{"deal_expired": 1}

		// Set user shopping list item deal expired to 1(true)
		err = ds.ShoppingListItemFactory.UpdateByUserGUIDAndDealGUID(userGUID, dealGUID, data)

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

	dealWithRelations := ds.DealRepository.GetDealByIDWithRelations(deal.ID, relations)

	return dealWithRelations
}

// GetAllDealsBasedOnLatitudeAndLongitudeAndQuota function used to retrieve all deals within valid range 10KM
func (ds *DealService) GetAllDealsBasedOnLatitudeAndLongitudeAndQuota(latitude string, longitude string, offset string, limit string, relations string) ([]*Deal, int) {
	// Convert Latitude and Longitude from string to float65
	latitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(latitude), 64)
	longitude1InFLoat64, _ := strconv.ParseFloat(strings.TrimSpace(longitude), 64)

	validDeals, totalDeal := ds.DealRepository.GetAllDealsWithinValidRangeAndQuota(latitude1InFLoat64, longitude1InFLoat64, offset, limit, relations)

	for key, validDeal := range validDeals {
		validDeals[key].Items = ds.ItemRepository.GetByID(validDeal.ItemID, "Categories,Subcategories")
	}

	return validDeals, totalDeal
}
