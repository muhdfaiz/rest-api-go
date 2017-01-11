package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// ItemCategoryTransformerInterface is a contract that defines the method needed for Item Category Transformer.
type ItemCategoryTransformerInterface interface {
	TransformCollection(data interface{}, totalData int) *ItemCategoryResponse
}

// ItemCategoryServiceInterface is a contract that defines the method needed for ItemCategoryService.
type ItemCategoryServiceInterface interface {
	GetItemCategoryNames() ([]string, int)
	GetItemCategories() ([]*ItemCategory, int)
	GetItemCategoryByGUID(guid string) *ItemCategory
	GetItemCategoryByID(itemCategoryID int) *ItemCategory
	TransformItemCategories(data interface{}, totalData int) *ItemCategoryResponse
	GetGrocerCategoriesThoseHaveDealsIncludingDeals(userGUID, grocerGUID,
		latitude, longitude, dealLimitPerCategory, relations string) ([]*ItemCategory, *systems.ErrorData)
}

// ItemCategoryRepositoryInterface is a function that defines the method needed for Item Category Repository.
type ItemCategoryRepositoryInterface interface {
	GetAll() ([]*ItemCategory, int)
	GetAllCategoryNames() ([]string, int)
	GetByID(ID int) *ItemCategory
	GetByGUID(GUID string) *ItemCategory
	GetGrocerCategoriesForThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID string, grocerID int,
		currentDateInGMT8 string, latitude, longitude float64) []*ItemCategory
}
