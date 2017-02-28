package v1_1

type NotificationServiceInterface interface {
	GetAllNotificationsForDevice(deviceUUID string) []*Notification
	GetNewsAndDealNotificationsForDevice(deviceUUID string) []*Notification
}

type NotificationRepositoryInterface interface {
	GetByDeviceUUID(deviceUUID, relations string) []*Notification
	GetByDeviceUUIDAndTypes(deviceUUID, relations string) []*Notification
}
