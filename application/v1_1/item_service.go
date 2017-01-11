package v1_1

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
