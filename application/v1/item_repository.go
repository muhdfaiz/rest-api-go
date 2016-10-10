package v1

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ItemRepositoryInterface interface {
	GetAll() []*Item
	GetLatestUpdate(lastSyncDate string) []*Item
}

type ItemRepository struct {
	DB *gorm.DB
}

func (or *ItemRepository) GetAll() []*Item {
	items := []*Item{}

	or.DB.Model(&Item{}).Find(&items)

	return items
}

func (or *ItemRepository) GetLatestUpdate(lastSyncDate string) []*Item {
	lastSync, _ := time.Parse(time.RFC3339, lastSyncDate)

	items := []*Item{}

	or.DB.Table("item").Where("updated_at > ?", lastSync).Order("updated_at desc").Find(&items)

	return items
}
