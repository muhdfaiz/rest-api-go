package v1

type ItemSubCategoryService struct {
	ItemSubCategoryRepository ItemSubCategoryRepositoryInterface
	DealRepository            DealRepositoryInterface
}

// GetItemSubCategoryByGUID function used to retrieve item subcategory by GUID.
func (iscs *ItemSubCategoryService) GetItemSubCategoryByGUID(itemSubCategoryGUID string) *ItemSubCategory {
	itemSubCategory := iscs.ItemSubCategoryRepository.GetByGUID(itemSubCategoryGUID)

	return itemSubCategory
}

// GetItemSubCategoryByID function used to retrieve item subcategory by ID.
func (iscs *ItemSubCategoryService) GetItemSubCategoryByID(itemSubCategoryID int) *ItemSubCategory {
	itemSubCategory := iscs.ItemSubCategoryRepository.GetByID(itemSubCategoryID)

	return itemSubCategory
}
