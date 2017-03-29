package v11

// ItemSubCategoryServiceInterface is a contract that defines the methods needed for Item Subcategory Service.
type ItemSubCategoryServiceInterface interface {
	GetItemSubCategoryByGUID(itemSubCategoryGUID string) *ItemSubCategory
	GetItemSubCategoryByID(itemSubCategoryID int) *ItemSubCategory
}

// ItemSubCategoryRepositoryInterface is a contract that defines the method needed for Item Subcategory Repository.
type ItemSubCategoryRepositoryInterface interface {
	GetByID(id int) *ItemSubCategory
	GetByGUID(guid string) *ItemSubCategory
	GetSubCategoriesForCategoryThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, categoryGUID, currentDateInGMT8 string, latitude,
		longitude float64) []*ItemSubCategory
}
