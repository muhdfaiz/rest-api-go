package v1

// DefaultShoppingListItemServiceInterface is a contract that defines the method needed
// for Default Shopping List Item Service.
type DefaultShoppingListItemServiceInterface interface {
	GetAllDefaultShoppingListItems() []*DefaultShoppingListItem
}

// DefaultShoppingListItemRepositoryInterface is a contract that defines the methods needed for
// Default Shopping List Item Repository.
type DefaultShoppingListItemRepositoryInterface interface {
	GetAll() []*DefaultShoppingListItem
}
