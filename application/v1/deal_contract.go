package v1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// DealTransformerInterface is a contract that defines the method needed for Deal Transformer.
type DealTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *DealResponse
}

// DealServiceInterface is a contract the defines the method needed for Deal Service.
type DealServiceInterface interface {
	CheckDealExistOrNotByGUID(dealGUID string) (*Deal, *systems.ErrorData)
	GetDealsBasedOnUserShoppingListItem(userGUID, shoppingListGUID string, shopppingListItems *ShoppingListItem, latitude,
		longitude string) []*Deal
	GetDealsBasedOnSampleShoppingListItem(defaultShoppingListItem *DefaultShoppingListItem, latitude,
		longitude string) []*Deal
	FilteredDealMustBeUniquePerShoppingList(deals []*Deal, dealsCollection []*Deal, userGUID string) []*Deal
	FilteredDealMustBeWithinStartAndEndTime(deals []*Deal, currentDateInGMT8, currentTimeInGMT8 string) []*Deal
	FilteredDealByPositiveTag(deals []*Deal, shoppingListItemName string) []*Deal
	FilteredDealsNotAddedTolist(deals []*Deal, userGUID, shoppingListGUID string) []*Deal
	GetFirstThreeDeals(deals []*Deal) []*Deal
	ViewDealDetails(dealGUID, relations string) *Ads
	GetAvailableDealsForGuestUser(latitude, longitude, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetAvailableDealsForRegisteredUser(userGUID, name, latitude, longitude, pageNumber, pageLimit,
		relations string) ([]*Deal, int)
	GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID, latitude, longitude, dealLimitPerCategory,
		relations string) []*ItemCategory
	GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID, categoryGUID, latitude, longitude,
		dealLimitPerSubcategory, relations string) []*ItemSubCategory
	GetAvailableDealsByCategoryForRegisteredUser(userGUID, category, latitude, longitude, pageNumber, pageLimit,
		relations string) ([]*Deal, int)
	GetAvailableDealsForSubCategoryForRegisteredUser(userGUID, category, latitude, longitude, pageNumber, pageLimit,
		relations string) ([]*Deal, int)
	GetAvailableDealsForGrocerByCategory(request *http.Request, userGUID, grocerGUID, categoryGUID, latitude, longitude,
		pageNumber, pageLimit, relations string) (*DealResponse, *systems.ErrorData)
	SetAddTolistInfoAndItemsAndGrocerExclusiveForDeal(deal *Deal, userGUID string) *Deal
	SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(deals []*Deal, userGUID string) []*Deal
}

// DealRepositoryInterface is a contract that defines the method needed for Deal Repository.
type DealRepositoryInterface interface {
	SumCashbackAmount(dealGUIDs []string) float64
	GetDealsByCategoryAndValidStartEndDate(todayDateInGMT8 string, shoppingListItem *ShoppingListItem) []*Deal
	GetDealsByValidStartEndDate(todayDateInGMT8 string) []*Deal
	GetDealByGUID(dealGUID string) *Deal
	GetDealByIDWithRelations(dealID int, relations string) *Ads
	GetDealsWithinRangeAndDateRangeAndQuota(latitude, longitude float64, currentDateInGMT8,
		pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsWithinRangeAndDateRangeAndUserLimitAndQuotaAndName(userGUID, name string, latitude, longitude float64,
		currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsForCategoryWithinDateRangeAndQuota(category, currentDateInGMT8,
		pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsForCategoryWithinRangeAndDateRangeAndQuota(category string, latitude, longitude float64,
		currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsByCategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, category string, latitude, longitude float64,
		currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsBySubcategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, subcategory string, latitude, longitude float64,
		currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuotaAndCategory(userGUID, categoryGUID string, grocerID int,
		latitude, longitude float64, currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	CountDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID string, grocerID int,
		latitude, longitude float64, currentDateInGMT8 string) int
	GetDealBySubCategoryGUIDWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, subCategoryGUID string, latitude, longitude float64,
		currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int)
	GetDealByGUIDAndValidStartEndDate(dealGUID, todayDateInGMT8 string) *Deal
}
