package v1

// ItemServiceInterface is a contract that defines the methods needed for Item Service
type ItemServiceInterface interface {
	GetItemByName(itemName, relations string) *Item
	GetItemByID(itemID int, relations string) *Item
}

type ItemService struct {
	ItemRepository ItemRepositoryInterface
}

// GetItemByName function used to retrieve Item by name
func (is *ItemService) GetItemByName(itemName, relations string) *Item {
	item := is.ItemRepository.GetByName(itemName, relations)

	return item
}

// GetItemByID function used to retrieve Item by Item GUID.
func (is *ItemService) GetItemByID(itemID int, relations string) *Item {
	item := is.ItemRepository.GetByID(itemID, relations)

	return item
}
