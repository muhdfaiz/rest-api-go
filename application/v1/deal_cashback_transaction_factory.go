package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type DealCashbackTransactionFactoryInterface interface {
	Create(userGUID string, transactionGUID, receiptURL string) (*DealCashbackTransaction, *systems.ErrorData)
}

type DealCashbackTransactionFactory struct {
	DB *gorm.DB
}

func (dctf *DealCashbackTransactionFactory) Create(userGUID string, transactionGUID, receiptURL string) (*DealCashbackTransaction, *systems.ErrorData) {
	createdDealCashbackTransaction := &DealCashbackTransaction{}

	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:            Helper.GenerateUUID(),
		UserGUID:        userGUID,
		TransactionGUID: transactionGUID,
		ReferenceID:     Helper.GenerateUniqueShortID(),
		ReceiptURL:      receiptURL,
	}

	result := dctf.DB.Create(dealCashbackTransaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdDealCashbackTransaction = result.Value.(*DealCashbackTransaction)

	return createdDealCashbackTransaction, nil
}
