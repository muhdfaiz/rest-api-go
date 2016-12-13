package v1

// ItemCategoryServiceInterface is a contract that defines the method needed for ItemCategoryService.
type ItemCategoryServiceInterface interface {
	GetItemCategoryNames() ([]string, int)
	GetItemCategories() ([]*ItemCategory, int)
	GetItemCategoryByGUID(guid string) *ItemCategory
	GetItemCategoryByID(itemCategoryID int) *ItemCategory
	TransformItemCategories(data interface{}, totalData int) *ItemCategoryResponse
}

type ItemCategoryService struct {
	ItemCategoryRepository  ItemCategoryRepositoryInterface
	ItemCategoryTransformer ItemCategoryTransformerInterface
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
