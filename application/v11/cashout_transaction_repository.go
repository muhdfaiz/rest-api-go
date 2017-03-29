package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// CashoutTransactionRepository will handle all CRUD function for Cashout Transaction resource.
type CashoutTransactionRepository struct {
	BaseRepository
	DB *gorm.DB
}

// Create new Cashout Transaction and store in database.
func (ctr *CashoutTransactionRepository) Create(dbTransaction *gorm.DB, userGUID, transactionGUID string,
	cashoutTransactionData *CreateCashoutTransaction) (*CashoutTransaction, *systems.ErrorData) {

	cashoutTransaction := &CashoutTransaction{
		GUID:                  Helper.GenerateUUID(),
		UserGUID:              userGUID,
		TransactionGUID:       transactionGUID,
		BankName:              cashoutTransactionData.BankName,
		BankAccountHolderName: cashoutTransactionData.BankAccountHolderName,
		BankAccountNumber:     cashoutTransactionData.BankAccountNumber,
		BankCountry:           cashoutTransactionData.BankCountry,
	}

	result := dbTransaction.Create(cashoutTransaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*CashoutTransaction), nil
}

// CountByUserGUID function used to count total number of cashout transaction filter by User GUID.
func (ctr *CashoutTransactionRepository) CountByUserGUID(userGUID string) int {
	var numberOfCashoutTransaction int

	ctr.DB.Model(&CashoutTransaction{}).Where(&CashoutTransaction{UserGUID: userGUID}).Count(&numberOfCashoutTransaction)

	return numberOfCashoutTransaction
}
