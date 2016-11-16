package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type DealCashbackTransactionFactoryInterface interface {
	Create(userGUID string, images map[string]string) (*DealCashbackTransaction, *systems.ErrorData)
}

type DealCashbackTransactionFactory struct {
	DB                                      *gorm.DB
	DealCashbackTransactionStatusRepository DealCashbackTransactionStatusRepositoryInterface
}

func (dctf *DealCashbackTransactionFactory) Create(userGUID string, image map[string]string) (*DealCashbackTransaction, *systems.ErrorData) {
	dealCashbackTransactionStatus := dctf.DealCashbackTransactionStatusRepository.GetBySlug("pending")

	createdDealCashbackTransaction := &DealCashbackTransaction{}

	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:                              Helper.GenerateUUID(),
		UserGUID:                          userGUID,
		ReceiptID:                         Helper.GenerateUniqueShortID(),
		ReceiptImage:                      image["path"],
		DealCashbackTransactionStatusGUID: dealCashbackTransactionStatus.GUID,
	}

	result := dctf.DB.Create(dealCashbackTransaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdDealCashbackTransaction = result.Value.(*DealCashbackTransaction)

	return createdDealCashbackTransaction, nil
}
