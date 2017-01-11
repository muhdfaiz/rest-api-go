package v1

import "net/http"

// ItemTransformerInterface is a contact that defines the method needed for Item Transformer.
type ItemTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *ItemResponse
}

// ItemServiceInterface is a contract that defines the methods needed for Item Service
type ItemServiceInterface interface {
	GetItemByName(itemName, relations string) *Item
	GetItemByID(itemID int, relations string) *Item
}

// ItemRepositoryInterface is a contract that defines the methods needed for Item Repository.
type ItemRepositoryInterface interface {
	GetAll(pageNumber string, pageLimit string, relations string) ([]*Item, int)
	GetLatestUpdate(lastSyncDate string, pageNumber string, pageLimit string, relations string) ([]*Item, int)
	GetByID(id int, relations string) *Item
	GetByName(name string, relations string) *Item
	GetUniqueCategories(relations string) ([]string, int)
}
