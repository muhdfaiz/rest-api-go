package v1_1

type NotificationService struct {
	NotificationRepository NotificationRepositoryInterface
}

func (ns *NotificationService) GetAllNotificationsForDevice(deviceUUID string) []*Notification {
	notifications := ns.NotificationRepository.GetByDeviceUUID(deviceUUID, "Transactions,Transactions.Transactiontypes,Transactions.Transactionstatuses")

	return notifications
}

func (ns *NotificationService) GetNewsAndDealNotificationsForDevice(deviceUUID string) []*Notification {
	notifications := ns.NotificationRepository.GetByDeviceUUIDAndTypes(deviceUUID, "Transactions,Transactions.Transactiontypes,Transactions.Transactionstatuses")

	return notifications
}
