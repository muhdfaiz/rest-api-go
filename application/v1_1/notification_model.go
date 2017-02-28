package v1_1

import "time"

type Notification struct {
	ID               int       `json:"id"`
	GUID             string    `json:"guid"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	Type             string    `json:"type"`
	Timestamp        string    `json:"timestamp"`
	ImageURL         string    `json:"image_url"`
	TransactionGUID  string    `json:"transaction_guid"`
	ReadNotification int       `json:"read_notification"`
	UUID             string    `json:"uuid"`
	AdminID          int       `json:"admin_id"`
	CreatedAt        time.Time `json:"created_at"`

	// Notification Has One Transaction
	Transactions *Transaction `json:"transaction,omitempty" gorm:"ForeignKey:TransactionGUID;AssociationForeignKey:GUID"`
}

// TableName function used to override default plural table name used by gorm based on struct name.
func (n Notification) TableName() string {
	return "notification"
}
