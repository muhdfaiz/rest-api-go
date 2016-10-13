package v1

import "github.com/jinzhu/gorm"

type ItemRepositoryInterface interface {
	GetAll(pageNumber string, pageLimit string) ([]*Item, int)
	GetLatestUpdate(lastSyncDate string, pageNumber string, pageLimit string) ([]*Item, int)
}

// ItemRepository will handle task related to retrieve and search shopping list items in database
type ItemRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all shopping list item insert from admin control panel
func (ir *ItemRepository) GetAll(pageNumber string, pageLimit string) ([]*Item, int) {
	items := []*Item{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	ir.DB.Model(&Item{}).Offset(offset).Limit(pageLimit).Order("updated_at desc").Find(&items)

	var totalItem *int

	ir.DB.Model(&Item{}).Count(&totalItem)

	return items, *totalItem
}

// GetLatestUpdate function used to retrieve latest update shopping list item that happen after last sync date in the query string
func (ir *ItemRepository) GetLatestUpdate(lastSyncDate string, pageNumber string, pageLimit string) ([]*Item, int) {
	items := []*Item{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	ir.DB.Model(&Item{}).Offset(offset).Limit(pageLimit).Where("updated_at > ?", lastSyncDate).Order("updated_at desc").Find(&items)

	var totalItem *int

	ir.DB.Model(&Item{}).Where("updated_at > ?", lastSyncDate).Count(&totalItem)

	return items, *totalItem
}
