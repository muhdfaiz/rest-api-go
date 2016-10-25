package v1

import "github.com/jinzhu/gorm"

type ShoppingListItemRepositoryInterface interface {
	GetByName(name string, relations string) *ShoppingListItem
	GetByGUID(guid string, relations string) *ShoppingListItem
	GetByShoppingListGUIDAndGUID(guid string, shoppingListGUID string, relations string) *ShoppingListItem
	GetByUserGUIDAndShoppingListGUID(userGUID string, shoppingListGUID string, relations string) []*ShoppingListItem
	GetUserShoppingListItem(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem
	GetUserShoppingListItemNotAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem
	GetUserShoppingListItemAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem
}

// ShoppingListItemRepository used to handle all task related to viewing, retrieving shopping list item
type ShoppingListItemRepository struct {
	DB *gorm.DB
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
func (slir *ShoppingListItemRepository) GetUserShoppingListItem(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem {
	shoppingListItemsGroupByCategory := []*ShoppingListItem{}

	userShoppingListItemsGroupByCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID}).Group("category").Find(&shoppingListItemsGroupByCategory)

	for _, shoppingListItemGroupByCategory := range shoppingListItemsGroupByCategory {
		userShoppingListItems := []*ShoppingListItem{}

		DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, Category: shoppingListItemGroupByCategory.Category}).Find(&userShoppingListItems)

		userShoppingListItemsGroupByCategory[shoppingListItemGroupByCategory.Category] = userShoppingListItems
	}

	return userShoppingListItemsGroupByCategory
}

// GetUserShoppingListItemAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that added to cart
func (slir *ShoppingListItemRepository) GetUserShoppingListItemAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem {
	shoppingListItemsGroupByCategory := []*ShoppingListItem{}

	userShoppingListItemsGroupByCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, AddedToCart: 1}).Group("category").
		Find(&shoppingListItemsGroupByCategory)

	for _, shoppingListItemGroupByCategory := range shoppingListItemsGroupByCategory {
		userShoppingListItems := []*ShoppingListItem{}

		DB.Where(&ShoppingListItem{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, Category: shoppingListItemGroupByCategory.Category, AddedToCart: 1}).
			Find(&userShoppingListItems)

		userShoppingListItemsGroupByCategory[shoppingListItemGroupByCategory.Category] = userShoppingListItems
	}

	return userShoppingListItemsGroupByCategory
}

// GetUserShoppingListItemNotAddedToCart function used to retrieve shopping list item by user guid and shopping list guid that not added to cart
func (slir *ShoppingListItemRepository) GetUserShoppingListItemNotAddedToCart(userGUID string, shoppingListGUID string, relations string) map[string][]*ShoppingListItem {
	shoppingListItemsGroupByCategory := []*ShoppingListItem{}

	userShoppingListItemsGroupByCategory := make(map[string][]*ShoppingListItem)

	DB := slir.DB.Model(&ShoppingListItem{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart != ?", userGUID, shoppingListGUID, 1).Group("category").
		Find(&shoppingListItemsGroupByCategory)

	for _, shoppingListItemGroupByCategory := range shoppingListItemsGroupByCategory {
		userShoppingListItems := []*ShoppingListItem{}

		DB.Where("user_guid = ? AND shopping_list_guid = ? AND added_to_cart != ? AND category = ?", userGUID, shoppingListGUID, 1, shoppingListItemGroupByCategory.Category).
			Find(&userShoppingListItems)

		userShoppingListItemsGroupByCategory[shoppingListItemGroupByCategory.Category] = userShoppingListItems
	}

	return userShoppingListItemsGroupByCategory
}
