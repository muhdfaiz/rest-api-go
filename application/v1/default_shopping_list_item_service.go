package v1

type DefaultShoppingListItemService struct {
	DefaultShoppingListItemRepository DefaultShoppingListItemRepositoryInterface
}

// GetAllDefaultShoppingListItems function used to retrieve all default shopping list items
// those were set in Admin Dashboard.
func (dslis *DefaultShoppingListItemService) GetAllDefaultShoppingListItems() []*DefaultShoppingListItem {
	defaultShoppingListItems := dslis.DefaultShoppingListItemRepository.GetAll()

	return defaultShoppingListItems
}
