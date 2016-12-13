package v1

// ItemServiceInterface is a contract that defines the methods needed for Item Service
type ItemServiceInterface interface {
	GetItemByName(itemName string) *Item
}

type ItemService struct {
	ItemRepository ItemRepositoryInterface
}

// GetItemByName function used to retrieve Item by name
func (is *ItemService) GetItemByName(itemName string) *Item {
	item := is.ItemRepository.GetByName(itemName, "")

	return item
}
