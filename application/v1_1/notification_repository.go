package v1_1

import "github.com/jinzhu/gorm"

// NotificationRepository will handle all CRUD operation related to Notification resource.
type NotificationRepository struct {
	DB *gorm.DB
}

// GetByDeviceUUIDAndBlastTypeAndEmptyUserGUIDAndType function used to retrive notification filter by device UUID and blastType
// type and empty user GUID and types.
func (nr *NotificationRepository) GetByDeviceUUIDAndBlastTypeAndEmptyUserGUIDAndType(deviceUUID, blastType, relations string) []*Notification {
	notifications := []*Notification{}

	DB := nr.DB.Model(&Notification{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Notification{UUID: deviceUUID, Blast: blastType}).Where("type = ? OR type = ?", "news", "deals").
		Where("user_guid IS NULL OR user_guid = ''").Find(&notifications)

	return notifications
}

func (nr *NotificationRepository) GetByUserGUIDOrUserGUIDEmptyAndDeviceUUID(deviceUUID, userGUID, blastType, relations string) []*Notification {
	notifications := []*Notification{}

	DB := nr.DB.Model(&Notification{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where("user_guid = '"+userGUID+"' OR user_guid IS NULL OR user_guid = ''").Where("uuid = ?", deviceUUID).Find(&notifications)

	return notifications
}
