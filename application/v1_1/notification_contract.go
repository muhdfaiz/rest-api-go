package v1_1

type NotificationServiceInterface interface {
	GetNotificationsForGuest(deviceUUID string) []*Notification
	GetNotificationsForLoggedInUser(deviceUUID, userGUID string) []*Notification
}

type NotificationRepositoryInterface interface {
	GetByDeviceUUIDAndBlastTypeAndEmptyUserGUIDAndType(deviceUUID, blastType, relations string) []*Notification
	GetByUserGUIDOrUserGUIDEmptyAndDeviceUUID(deviceUUID, userGUID, blastType, relations string) []*Notification
}
