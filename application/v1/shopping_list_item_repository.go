package v1

import "github.com/jinzhu/gorm"

type ShoppingListItemRepositoryInterface interface {
	GetByName(name string, relations string) *ShoppingListItem
	GetByGUID(guid string, relations string) *ShoppingListItem
	GetByShoppingListGUIDAndGUID(guid string, shoppingListGUID string, relations string) *ShoppingListItem
	GetByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, relations string) []*ShoppingListItem
	GetUserShoppingListItem(userGUID string, shoppingListGUID string, relations string, latitude string,
		longitude string) map[string][]*ShoppingListItem
	GetUserShoppingListItemAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem
	GetUserShoppingListItemNotAddedToCart(userGUID string, shoppingListGUID string, relations string, latitude string,
		longitude string) map[string][]*ShoppingListItem
}

// ShoppingListItemRepository used to handle all task related to viewing, retrieving shopping list item
type ShoppingListItemRepository struct {
	DB          *gorm.DB
	DealService DealServiceInterface
}

// GetByGUID function used to retrieve shopping list item by GUID
func (slir *ShoppingListItemRepository) GetByGUID(guid string, relations string) *ShoppingListItem {

	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{GUID: guid}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByName function used to retrieve shopping list item by Name
func (slir *ShoppingListItemRepository) GetByName(name string, relations string) *ShoppingListItem {
	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{Name: name}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByShoppingListGUIDAndGUID function used to retrieve shopping list item by shopping list GUID and GUID
func (slir *ShoppingListItemRepository) GetByShoppingListGUIDAndGUID(guid string, shoppingListGUID string, relations string) *ShoppingListItem {
	shoppingListItem := &ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{GUID: guid, ShoppingListGUID: shoppingListGUID}).First(&shoppingListItem)

	return shoppingListItem
}

// GetByUserGUIDAndShoppingListGUID function used to retrieve shopping list item by user GUID and shopping list GUID
func (slir *ShoppingListItemRepository) GetByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, relations string) []*ShoppingListItem {
	shoppingListItem := []*ShoppingListItem{}

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Find(&shoppingListItem)

	return shoppingListItem
}

// GetUserShoppingListItem function used to retrieve shopping list item by user GUID and shopping list GUID
func (slir *ShoppingListItemRepository) GetUserShoppingListItem(userGUID string, shoppingListGUID string, relations string, latitude string,
	longitude string) map[string][]*ShoppingListItem {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	shoppingListItemsGroupBySubCategory := []*ShoppingListItem{}

	dealsCollection := []*Deal{}

	// Retrieve unique shopping list item category from user shopping list
	// Get Unique Category for user shopping list items
	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Group("sub_category").
		Find(&shoppingListItemsGroupBySubCategory)

	// Loop through each shopping list item category
	for _, shoppingListItemGroupBySubCategory := range shoppingListItemsGroupBySubCategory {
		userShoppingListItems := []*ShoppingListItem{}

		// Retrieve shopping list item by shopping list item category
		DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, SubCategory: shoppingListItemGroupBySubCategory.SubCategory}).Find(&userShoppingListItems)

		// Retrieve available deals for each item. Maximum deal per item is 3
		for key, userShopppingListItem := range userShoppingListItems {

			// If user shopping list item was not added from deal and not added to cart, retrieve valid deals
			if userShopppingListItem.AddedFromDeal == 0 && userShopppingListItem.AddedToCart == 0 && latitude != "" && longitude != "" {
				deals := slir.DealService.GetDealsBasedOnUserShoppingListItem(userGUID, userShopppingListItem, latitude, longitude, dealsCollection)

				userShoppingListItems[key].Deals = nil

				if len(deals) > 0 {
					dealsCollection = append(dealsCollection, deals...)
					userShoppingListItems[key].Deals = deals
				}
			}

			// If user shopping list item was added from deal and not added to cart, check deal expired or not
			if userShopppingListItem.AddedFromDeal == 1 && userShopppingListItem.AddedToCart == 0 {
				slir.DealService.RemoveDealCashbackAndSetItemDealExpired(userGUID, *userShopppingListItem.DealGUID)
			}
		}

		userShoppingListItemsGroupBySubCategory[shoppingListItemGroupBySubCategory.SubCategory] = userShoppingListItems
	}

	return userShoppingListItemsGroupBySubCategory
}

// GetUserShoppingListItemAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that added to cart
func (slir *ShoppingListItemRepository) GetUserShoppingListItemAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem {
	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	shoppingListItemsGroupBySubCategory := []*ShoppingListItem{}

	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 1}).Group("sub_category").
		Find(&shoppingListItemsGroupBySubCategory)

	for _, shoppingListItemGroupBySubCategory := range shoppingListItemsGroupBySubCategory {
		userShoppingListItems := []*ShoppingListItem{}

		DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, SubCategory: shoppingListItemGroupBySubCategory.SubCategory, AddedToCart: 1}).
			Find(&userShoppingListItems)

		userShoppingListItemsGroupBySubCategory[shoppingListItemGroupBySubCategory.SubCategory] = userShoppingListItems
	}

	return userShoppingListItemsGroupBySubCategory
}

// GetUserShoppingListItemNotAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that not added to cart
func (slir *ShoppingListItemRepository) GetUserShoppingListItemNotAddedToCart(userGUID string, shoppingListGUID string, relations string,
	latitude string, longitude string) map[string][]*ShoppingListItem {

	userShoppingListItemsGroupBySubCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	shoppingListItemsGroupBySubCategory := []*ShoppingListItem{}

	// Get Unique Category for user shopping list items
	DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart != ?", userGUID, shoppingListGUID, 1).Group("sub_category").
		Find(&shoppingListItemsGroupBySubCategory)

	dealsCollection := []*Deal{}

	// Loop through each of user shopping list item category
	for _, shoppingListItemGroupBySubCategory := range shoppingListItemsGroupBySubCategory {
		userShoppingListItems := []*ShoppingListItem{}

		DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart != ? AND sub_category = ?", userGUID, shoppingListGUID, 1, shoppingListItemGroupBySubCategory.SubCategory).
			Find(&userShoppingListItems)

		// Retrieve available deals for each item. Maximum deal per item is 3
		for key, userShopppingListItem := range userShoppingListItems {

			// If user shopping list item was not added from deal and not added to cart, retrieve valid deals
			if userShopppingListItem.AddedFromDeal == 0 && userShopppingListItem.AddedToCart == 0 && latitude != "" && longitude != "" {
				deals := slir.DealService.GetDealsBasedOnUserShoppingListItem(userGUID, userShopppingListItem, latitude, longitude, dealsCollection)

				dealsCollection = append(dealsCollection, deals...)

				userShoppingListItems[key].Deals = deals
			}

			// If user shopping list item was added from deal and not added to cart, check deal expired or not
			if userShopppingListItem.AddedFromDeal == 1 && userShopppingListItem.AddedToCart == 0 {
				slir.DealService.RemoveDealCashbackAndSetItemDealExpired(userGUID, *userShopppingListItem.DealGUID)
			}
		}

		userShoppingListItemsGroupBySubCategory[shoppingListItemGroupBySubCategory.SubCategory] = userShoppingListItems
	}

	return userShoppingListItemsGroupBySubCategory
}
