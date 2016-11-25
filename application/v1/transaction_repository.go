package v1

import (
	"strings"

	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// TransactionRepositoryInterface is a contract that define the methods needed for Transaction Repository
type TransactionRepositoryInterface interface {
	Create(createTransactionData *CreateTransaction) (*Transaction, *systems.ErrorData)
	GetByGUID(GUID string, relations string) *Transaction
	GetByUserGUIDAndStatusAndReadStatus(userGUID string, transactionStatus string, readStatus string, pageNumber string,
		pageLimit string, relations string) ([]*Transaction, int)
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

// GetByUserGUIDAndStatusAndReadStatus function used to retrieve transactions by User GUID and Transaction Status and Read Status.
// If the value Read Status is 1 means user already read or click the transaction.
// If the value is 0 means user still not read or click the transaction.
func (tr *TransactionRepository) GetByUserGUIDAndStatusAndReadStatus(userGUID string, transactionStatus string, readStatus string,
	pageNumber string, pageLimit string, relations string) ([]*Transaction, int) {

	transactions := []*Transaction{}

	DB := tr.DB.Model(&Transaction{})

	offset := SetOffsetValue(pageNumber, pageLimit)

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB = DB.Joins("LEFT JOIN transaction_statuses ON transaction_statuses.guid = transactions.transaction_status_guid").Where(&Transaction{UserGUID: userGUID})

	if transactionStatus != "" {
		transactionStatuses := strings.Split(transactionStatus, ",")

		for key, transactionStatus := range transactionStatuses {
			if key == 0 {
				DB = DB.Where("transaction_statuses.slug = ?", transactionStatus)
			} else {
				DB = DB.Or("transaction_statuses.slug = ?", transactionStatus)
			}
		}
	}

	if readStatus != "" {
		readStatusInInt, _ := strconv.Atoi(readStatus)
		DB = DB.Where(&Transaction{ReadStatus: readStatusInInt}, transactionStatus)
	}

	if pageLimit != "" && pageNumber != "" {
		DB.Offset(offset).Limit(pageLimit).Find(&transactions)
	} else {
		DB.Find(&transactions)
	}

	var TotalTransaction int

	DB.Count(&TotalTransaction)

	return transactions, TotalTransaction
}
