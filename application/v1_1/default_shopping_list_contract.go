package v1

// DefaultShoppingListServiceInterface is a contract that defines the method needed for
// Default Shopping List Service.
type DefaultShoppingListServiceInterface interface {
	GetAllDefaultShoppingLists(relations string) []*DefaultShoppingList
	GetAllDefaultShoppingListsIncludingItemsAndDeals(latitude string, longitude string,
		relations string) []*DefaultShoppingList
}

// DefaultShoppingListRepositoryInterface is a contract that defines the method needed for Default Shopping
// List Repository.
type DefaultShoppingListRepositoryInterface interface {
	GetAll(relations string) []*DefaultShoppingList
}
