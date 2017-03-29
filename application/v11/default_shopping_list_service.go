package v11

// DefaultShoppingListService used to handle application logic related to Default Shopping List resource.
type DefaultShoppingListService struct {
	DefaultShoppingListRepository DefaultShoppingListRepositoryInterface
	DealService                   DealServiceInterface
}

// GetAllDefaultShoppingLists function used to retrieve all default shopping List.
func (dsls *DefaultShoppingListService) GetAllDefaultShoppingLists(relations string) []*DefaultShoppingList {

	defaultShoppingLists := dsls.DefaultShoppingListRepository.GetAll(relations)

	return defaultShoppingLists
}

// GetAllDefaultShoppingListsIncludingItemsAndDeals function used to retrieve all default shopping List
// including items and deals.
func (dsls *DefaultShoppingListService) GetAllDefaultShoppingListsIncludingItemsAndDeals(latitude string, longitude string,
	relations string) []*DefaultShoppingList {

	defaultShoppingLists := dsls.DefaultShoppingListRepository.GetAll(relations)

	dealsCollection := []*Deal{}

	for key, defaultShoppingList := range defaultShoppingLists {
		for key1, defaultShoppingListItem := range defaultShoppingList.Items {

			deals := dsls.DealService.GetDealsBasedOnSampleShoppingListItem(defaultShoppingListItem, latitude, longitude)

			deals = dsls.DealService.FilteredDealMustBeUniquePerShoppingList(deals, dealsCollection, "")

			dealsCollection = append(dealsCollection, deals...)

			defaultShoppingLists[key].Items[key1].Deals = deals
		}
	}

	return defaultShoppingLists
}
