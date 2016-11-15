package v1

type ItemCategoryServiceInterface interface {
	GetItemCategoryNames() ([]string, int)
	GetItemCategories() ([]*ItemCategory, int)
	GetItemCategoryByGUID(guid string) *ItemCategory
	TransformItemCategories(data interface{}, totalData int) *ItemCategoryResponse
}

type ItemCategoryService struct {
	ItemCategoryRepository  ItemCategoryRepositoryInterface
	ItemCategoryTransformer ItemCategoryTransformerInterface
}

func (ics *ItemCategoryService) GetItemCategoryNames() ([]string, int) {
	itemCategories, totalItemCategory := ics.ItemCategoryRepository.GetAllCategoryNames()

	return itemCategories, totalItemCategory
}

func (ics *ItemCategoryService) GetItemCategories() ([]*ItemCategory, int) {
	itemCategories, totalItemCategory := ics.ItemCategoryRepository.GetAll()

	return itemCategories, totalItemCategory
}

func (ics *ItemCategoryService) GetItemCategoryByGUID(guid string) *ItemCategory {
	return ics.ItemCategoryRepository.GetByGUID(guid)
}

func (ics *ItemCategoryService) TransformItemCategories(data interface{}, totalData int) *ItemCategoryResponse {
	return ics.ItemCategoryTransformer.TransformCollection(data, totalData)
}
