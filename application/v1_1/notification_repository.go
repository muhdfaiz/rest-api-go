package v1_1

import "github.com/jinzhu/gorm"

// NotificationRepository will handle all CRUD operation related to Notification resource.
type NotificationRepository struct {
	DB *gorm.DB
}

// GetByDeviceUUID function used to retrive notification filter by device UUID.
func (nr *NotificationRepository) GetByDeviceUUID(deviceUUID, relations string) []*Notification {
	notifications := []*Notification{}

	DB := nr.DB.Model(&Notification{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Notification{UUID: deviceUUID}).Find(&notifications)

	return notifications
}

// GetByDeviceUUIDAndTypes function used to retrive notification filter by device UUID and notification types.
func (nr *NotificationRepository) GetByDeviceUUIDAndTypes(deviceUUID, relations string) []*Notification {
	notifications := []*Notification{}

	DB := nr.DB.Model(&Notification{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Notification{UUID: deviceUUID}).Where("type = ? OR type = ?", "news", "deals").Find(&notifications)

	return notifications
}
