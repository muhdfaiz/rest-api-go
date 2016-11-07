package v1

type ItemCategoryServiceInterface interface {
	GetAllItemCategoryNames() *ItemCategoryResponse
}

type ItemCategoryService struct {
	ItemCategoryRepository  ItemCategoryRepositoryInterface
	ItemCategoryTransformer ItemCategoryTransformerInterface
}

func (ics *ItemCategoryService) GetAllItemCategoryNames() *ItemCategoryResponse {
	itemCategories, totalItemCategory := ics.ItemCategoryRepository.GetAllCategoryNames()

	return ics.ItemCategoryTransformer.TransformCollection(itemCategories, totalItemCategory)
}
