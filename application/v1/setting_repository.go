package v1

import "github.com/jinzhu/gorm"

// SettingRepository will handle all CRUD functions related to Setting resource.
type SettingRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all settings from database.
func (sr *SettingRepository) GetAll() []*Setting {
	settings := []*Setting{}

	sr.DB.Model(&Setting{}).Find(&settings)

	return settings
}

// GetBySlug function used to retrieve setting by slug from database.
func (sr *SettingRepository) GetBySlug(slug string) *Setting {
	setting := &Setting{}

	sr.DB.Model(&Setting{}).Where(&Setting{Slug: slug}).Find(&setting)

	return setting
}
