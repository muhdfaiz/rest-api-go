package v1

// ItemSubCategoryServiceInterface is a contract that defines the methods needed for ItemSubCategoryService.
type ItemSubCategoryServiceInterface interface {
	GetItemSubCategoryByGUID(itemSubCategoryGUID string) *ItemSubCategory
	GetItemSubCategoryByID(itemSubCategoryID int) *ItemSubCategory
}

type ItemSubCategoryService struct {
	ItemSubCategoryRepository ItemSubCategoryRepositoryInterface
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
