package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// DealCashbackTransactionRepository will handle all CRUD function for Deal Cashback Transaction resource.
type DealCashbackTransactionRepository struct {
	DB *gorm.DB
}

// Create function used to create new deal cashback transaction and store in database.
func (dctr *DealCashbackTransactionRepository) Create(dbTransaction *gorm.DB, userGUID string, transactionGUID, receiptURL string) (*DealCashbackTransaction, *systems.ErrorData) {
	createdDealCashbackTransaction := &DealCashbackTransaction{}

	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:            Helper.GenerateUUID(),
		UserGUID:        userGUID,
		TransactionGUID: transactionGUID,
		ReceiptURL:      receiptURL,
		RemarkTitle:     nil,
		RemarkBody:      nil,
		Status:          "pending cleaning",
	}

	result := dbTransaction.Create(dealCashbackTransaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdDealCashbackTransaction = result.Value.(*DealCashbackTransaction)

	return createdDealCashbackTransaction, nil
}

// CountByUserGUID function used to count total number of deal cashback transaction made by user via user GUID.
func (dctr *DealCashbackTransactionRepository) CountByUserGUID(userGUID string) *int {
	var totalDealCashback *int

	dctr.DB.Model(&DealCashbackTransaction{}).Where("user_guid = ?", userGUID).Count(&totalDealCashback)

	return totalDealCashback
}
