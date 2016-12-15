package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// TransactionRepositoryInterface is a contract that define the methods needed for Transaction Repository
type TransactionRepositoryInterface interface {
	Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData)
	UpdateReadStatus(transactionGUID string, readStatus int) *systems.ErrorData
	GetByGUID(GUID string, relations string) *Transaction
	GetByUserGUID(userGUID string, relations string) []*Transaction
	SumTotalAmountOfUserPendingTransactions(userGUID string) float64
	SumTotalAmountOfUserCashoutTransaction(userGUID string) float64
}

// TransactionRepository contains all function that can be used for CRUD operations.
type TransactionRepository struct {
	DB                          *gorm.DB
	TransactionStatusRepository TransactionStatusRepositoryInterface
}

// Create function used to create transaction and store in Database
func (tr *TransactionRepository) Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData) {
	pendingTransactionStatusGUID := tr.TransactionStatusRepository.GetBySlug("pending").GUID

	createdTransaction := &Transaction{}

	transaction := &Transaction{
		GUID:                  Helper.GenerateUUID(),
		ReferenceID:           Helper.GenerateUniqueShortID(),
		UserGUID:              createTransactionData.UserGUID,
		TransactionTypeGUID:   createTransactionData.TransactionTypeGUID,
		TransactionStatusGUID: pendingTransactionStatusGUID,
		TotalAmount:           createTransactionData.Amount,
		ApprovedAmount:        nil,
		RejectedAmount:        nil,
	}

	result := tr.DB.Create(transaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdTransaction = result.Value.(*Transaction)

	return createdTransaction, nil
}

// UpdateReadStatus function used to set `read_status` column to new value
func (tr *TransactionRepository) UpdateReadStatus(transactionGUID string, readStatus int) *systems.ErrorData {
	updateResult := tr.DB.Model(&Transaction{}).Where(&Transaction{GUID: transactionGUID}).Update("read_status", readStatus)

	if updateResult.Error != nil {
		return Error.InternalServerError(updateResult.Error, systems.DatabaseError)
	}

	return nil
}

// GetByGUID function used to retrieve transaction details
func (tr *TransactionRepository) GetByGUID(GUID string, relations string) *Transaction {
	transaction := &Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{GUID: GUID}).First(&transaction)

	return transaction
}

// GetByUserGUID function used to retrieve transactions by user GUID.
func (tr *TransactionRepository) GetByUserGUID(userGUID string, relations string) []*Transaction {

	transactions := []*Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{UserGUID: userGUID}).Find(&transactions)

	return transactions
}

// SumTotalAmountOfUserPendingTransactions function used to sum all of total amount for deal cashback transaction
// with status pending.
func (tr *TransactionRepository) SumTotalAmountOfUserPendingTransactions(userGUID string) float64 {

	type PendingDealCashbackTransaction struct {
		TotalAmountOfPendingDealCashbackTransaction float64 `json:"total_amount_of_pending_deal_cashback_transaction"`
	}

	pendingDealCashbackTransaction := &PendingDealCashbackTransaction{}

	tr.DB.Model(&Transaction{}).Select("sum(transactions.total_amount) as total_amount_of_pending_deal_cashback_transaction").
		Joins("left join transaction_statuses on transaction_statuses.guid = transactions.transaction_status_guid").
		Joins("left join transaction_types on transaction_types.guid = transactions.transaction_type_guid").
		Where(&Transaction{UserGUID: userGUID}).Where("transaction_types.slug = ?", "deal_redemption").
		Where("transaction_statuses.slug = ?", "pending").Scan(pendingDealCashbackTransaction)

	return pendingDealCashbackTransaction.TotalAmountOfPendingDealCashbackTransaction
}

// SumTotalAmountOfUserCashoutTransaction function used to sum total amount of cashout transaction made by user.
func (tr *TransactionRepository) SumTotalAmountOfUserCashoutTransaction(userGUID string) float64 {

	type CashoutTransaction struct {
		TotalAmountOfCashoutTransaction float64 `json:"total_amount_of_cashout_transaction"`
	}

	cashoutTransaction := &CashoutTransaction{}

	tr.DB.Model(&Transaction{}).Select("sum(transactions.total_amount) as total_amount_of_cashout_transaction").
		Joins("left join transaction_statuses on transaction_statuses.guid = transactions.transaction_status_guid").
		Joins("left join transaction_types on transaction_types.guid = transactions.transaction_type_guid").
		Where(&Transaction{UserGUID: userGUID}).Where("transaction_types.slug = ?", "cashout").
		Where("transaction_statuses.slug = ?", "approved").Scan(cashoutTransaction)

	return cashoutTransaction.TotalAmountOfCashoutTransaction
}
