package v1

import "github.com/jinzhu/gorm"

// ItemRepository will handle task related to retrieve and search shopping list items in database
type ItemRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all shopping list item insert from admin control panel
func (ir *ItemRepository) GetAll(pageNumber string, pageLimit string, relations string) ([]*Item, int) {
	items := []*Item{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	DB := ir.DB.Table("item").Select("item.*, category.name as category")

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Joins("LEFT JOIN category ON item.category_id = category.id").Offset(offset).Limit(pageLimit).Order("updated_at desc").Find(&items)

	var totalItem *int

	ir.DB.Model(&Item{}).Count(&totalItem)

	return items, *totalItem
}

// GetLatestUpdate function used to retrieve latest update shopping list item that happen after last sync date in the query string
func (ir *ItemRepository) GetLatestUpdate(lastSyncDate string, pageNumber string, pageLimit string, relations string) ([]*Item, int) {
	items := []*Item{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	DB := ir.DB.Table("item").Select("item.*, category.name as category")

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Joins("LEFT JOIN category ON item.category_id = category.id").Where("item.updated_at > ?", lastSyncDate).Offset(offset).Limit(pageLimit).Order("updated_at desc").Find(&items)

	var totalItem *int

	ir.DB.Model(&Item{}).Where("updated_at > ?", lastSyncDate).Count(&totalItem)

	return items, *totalItem
}

// GetByID function used to retrieve shopping list item by GUID
func (ir *ItemRepository) GetByID(id int, relations string) *Item {
	item := &Item{}

	DB := ir.DB.Model(&Item{}).Where(&Item{ID: id})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.First(&item)

	return item
}

// GetByName function used to retrieve shopping list item by name
func (ir *ItemRepository) GetByName(name string, relations string) *Item {
	item := &Item{}

	ir.DB.Model(&Item{}).Where(&Item{Name: name}).First(&item)

	return item
}

// GetUniqueCategories function used to retrieve unique shopping list item category
func (ir *ItemRepository) GetUniqueCategories(relations string) ([]string, int) {
	items := []*Item{}
	var shoppingListItemCategories []string

	ir.DB.Model(&Item{}).Group("category").Find(&items).Pluck("category", &shoppingListItemCategories)

	return shoppingListItemCategories, len(shoppingListItemCategories)
}
