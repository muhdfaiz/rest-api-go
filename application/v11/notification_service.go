package v11

// NotificationService used to handle application logic related to Notification resource.
type NotificationService struct {
	NotificationRepository NotificationRepositoryInterface
}

// GetNotificationsForGuest function used to retrieve notification guest including.
func (ns *NotificationService) GetNotificationsForGuest(deviceUUID string) []*Notification {
	notifications := ns.NotificationRepository.GetByDeviceUUIDAndBlastTypeAndEmptyUserGUIDAndType(deviceUUID, "all", "Transactions,Transactions.Transactiontypes,Transactions.Transactionstatuses")

	return notifications
}

func (ns *NotificationService) GetNotificationsForLoggedInUser(deviceUUID, userGUID string) []*Notification {
	notifications := ns.NotificationRepository.GetByUserGUIDOrUserGUIDEmptyAndDeviceUUID(deviceUUID, userGUID, "targeted", "Transactions,Transactions.Transactiontypes,Transactions.Transactionstatuses")

	return notifications
}
