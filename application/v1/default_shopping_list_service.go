package v1

// DefaultShoppingListServiceInterface is a contract that defines the method needed for
// Default Shopping List Service.
type DefaultShoppingListServiceInterface interface {
	GetAllDefaultShoppingList() []*DefaultShoppingList
}

type DefaultShoppingListService struct {
	DefaultShoppingListRepository DefaultShoppingListRepositoryInterface
}

// GetAllDefaultShoppingList function used to retrieve all default shopping List
// set from Admin Dashboard.
func (dsls *DefaultShoppingListService) GetAllDefaultShoppingList() []*DefaultShoppingList {
	defaultShoppingLists := dsls.DefaultShoppingListRepository.GetAll()

	return defaultShoppingLists
}
