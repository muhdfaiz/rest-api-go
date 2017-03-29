package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// TransactionRepository contains all function that can be used for CRUD operations.
type TransactionRepository struct {
	BaseRepository
	DB                          *gorm.DB
	TransactionStatusRepository TransactionStatusRepositoryInterface
}

// Create function used to create new transaction and store in database.
// It's return newly created transaction or error if encountered.
func (tr *TransactionRepository) Create(dbTransaction *gorm.DB, createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData) {
	createdTransaction := &Transaction{}

	transaction := &Transaction{
		GUID:                  Helper.GenerateUUID(),
		ReferenceID:           Helper.GenerateUniqueShortID(),
		UserGUID:              createTransactionData.UserGUID,
		TransactionTypeGUID:   createTransactionData.TransactionTypeGUID,
		TransactionStatusGUID: createTransactionData.TransactionStatusGUID,
		TotalAmount:           createTransactionData.Amount,
		ApprovedAmount:        nil,
		RejectedAmount:        nil,
	}

	result := dbTransaction.Create(transaction)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	createdTransaction = result.Value.(*Transaction)

	return createdTransaction, nil
}

// UpdateReadStatus function used to update `read_status` field in transaction table to new value using transaction GUID.
// Available value for `readStatus` parameter:
// - 0; means user still not read view the transaction.
// - 1; means user already read or view the transaction.
func (tr *TransactionRepository) UpdateReadStatus(dbTransaction *gorm.DB, transactionGUID string, readStatus int) *systems.ErrorData {
	updateResult := dbTransaction.Model(&Transaction{}).Where(&Transaction{GUID: transactionGUID}).Update("read_status", readStatus)

	if updateResult.Error != nil {
		return Error.InternalServerError(updateResult.Error, systems.DatabaseError)
	}

	return nil
}

// GetByGUID function used to retrieve first transaction details found filter by transaction GUID.
func (tr *TransactionRepository) GetByGUID(GUID string, relations string) *Transaction {
	transaction := &Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = tr.LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{GUID: GUID}).First(&transaction)

	return transaction
}

// GetByUserGUID function used to retrieve first transactions found filter by user GUID.
// Able to load transaction relation like transaction status and transaction type.
func (tr *TransactionRepository) GetByUserGUID(userGUID string, relations string) []*Transaction {
	transactions := []*Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = tr.LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{UserGUID: userGUID}).Find(&transactions)

	return transactions
}

// GetByUserGUIDAndTransactionTypeGUIDAndTransactionStatusGUID function used to retrieve multiple transactions filter by fields below:
// - filter by user GUID.
// - filter by transaction type GUID.
// - filter by transaction status GUID.
func (tr *TransactionRepository) GetByUserGUIDAndTransactionTypeGUIDAndTransactionStatusGUID(userGUID, transactionTypeGUID, transactionStatusGUID, relations string) []*Transaction {
	transactions := []*Transaction{}

	DB := tr.DB.Model(&Transaction{})

	if relations != "" {
		DB = tr.LoadRelations(DB, relations)
	}

	DB.Where(&Transaction{UserGUID: userGUID, TransactionTypeGUID: transactionTypeGUID, TransactionStatusGUID: transactionStatusGUID}).Find(&transactions)

	return transactions
}

// SumTotalAmountOfUserPendingTransactions function used to sum all of total amount for deal cashback transaction
// with status `pending` for specific user using user GUID.
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

// SumTotalAmountOfUserCashoutTransaction function used to sum total amount of cashout transaction for specific
// user using user GUID with status = `approved` and transaction_types.slug = `cashout`
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
