package v1

// DefaultShoppingListItemServiceInterface is a contract that defines the method needed
// for Default Shopping List Item Service.
type DefaultShoppingListItemServiceInterface interface {
	GetAllDefaultShoppingListItems() []*DefaultShoppingListItem
}

type DefaultShoppingListItemService struct {
	DefaultShoppingListItemRepository DefaultShoppingListItemRepositoryInterface
}

// GetAllDefaultShoppingListItems function used to retrieve all default shopping list items
// those were set in Admin Dashboard.
func (dslis *DefaultShoppingListItemService) GetAllDefaultShoppingListItems() []*DefaultShoppingListItem {
	defaultShoppingListItems := dslis.DefaultShoppingListItemRepository.GetAll()

	return defaultShoppingListItems
}
