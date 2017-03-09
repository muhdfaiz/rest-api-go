package v1_1

import (
	"fmt"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"
)

func TestCreateCashoutTransactionShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestCreateCashoutTransactionShouldReturnRequiredFieldValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwtToken, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, users[0].GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		BankAccountHolderName: "",
		BankAccountNumber:     "",
		BankCountry:           "",
		BankName:              "",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	errorDetail := errors["detail"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errorDetail["amount"])
	require.NotEmpty(t, errorDetail["bank_account_name"])
	require.NotEmpty(t, errorDetail["bank_account_number"])
	require.NotEmpty(t, errorDetail["bank_name"])
	require.NotEmpty(t, errorDetail["bank_country"])
}

func TestCreateCashoutTransactionShouldReturnAmountGreaterThanZeroValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwtToken, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, users[0].GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                -1,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "3123123",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	errorDetail := errors["detail"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errorDetail["amount"])
}

func TestCreateCashoutTransactionShouldReturnNumericValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwtToken, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, users[0].GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                20.00,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "asdadasdasd",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	errorDetail := errors["detail"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errorDetail["bank_account_number"])
}

func TestCreateCashoutTransactionShouldReturnCashoutAmountExceededLimitError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("612345678", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, user.GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                40.00,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "1234567890",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Cashout Amount Exceeded Limit.", errors["title"])
}

func TestCreateCashoutTransactionShouldErrorIfUserStillHasPendingCashoutTransaction(t *testing.T) {
	// Prepare sample data
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	sampleData.ShoppingListItems(users[0].GUID, shoppingList.GUID, 0)

	sampleData.Categories()

	sampleData.Subcategories()

	sampleData.Generics()

	sampleData.Items()

	deals := sampleData.Deals()

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Pending Transaction Status
	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Sample transaction, dealCashback
	dealCashback := sampleData.DealCashback(users[0].GUID, shoppingList.GUID, deals[0].GUID, nil)

	transaction := sampleData.Transaction(users[0].GUID, cashoutTransactionType.GUID, pendingTransactionStatus.GUID, 0, deals[0].CashbackAmount)

	sampleData.DealCashbackTransactionWithPendingCleaningStatus(dealCashback.GUID, users[0].GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, users[0].GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                30.00,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "1234567890",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Pending Cashout Transaction.", errors["title"])
}

func TestCreateCashoutTransactionShouldCreateTransactionRecord(t *testing.T) {
	// Prepare sample data
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("612345678", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(user.GUID, occasions[0].GUID, "Test Shopping List")

	sampleData.ShoppingListItems(user.GUID, shoppingList.GUID, 0)

	sampleData.Categories()

	sampleData.Subcategories()

	sampleData.Generics()

	sampleData.Items()

	deals := sampleData.Deals()

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Approved Transaction Status
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	// Retrieve GUID for Pending Transaction Status
	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Sample transaction, dealCashback
	dealCashback := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	transaction := sampleData.Transaction(user.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 30.00)

	sampleData.DealCashbackTransactionWithCompletedStatus(dealCashback.GUID, user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, user.GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                25.00,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "1234567890",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.NotEmpty(t, data["guid"])
	require.Equal(t, user.GUID, data["user_guid"])
	require.Equal(t, cashoutTransactionType.GUID, data["transaction_type_guid"])
	require.Equal(t, pendingTransactionStatus.GUID, data["transaction_status_guid"])
	require.Equal(t, 0, int(data["read_status"].(interface{}).(float64)))
	require.NotEmpty(t, data["reference_id"])
	require.Equal(t, cashoutTransactionData.Amount, data["total_amount"])
}

func TestCreateCashoutTransactionShouldCreateCashoutTransactionRecord(t *testing.T) {
	// Prepare sample data
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("612345678", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(user.GUID, occasions[0].GUID, "Test Shopping List")

	sampleData.ShoppingListItems(user.GUID, shoppingList.GUID, 0)

	sampleData.Categories()

	sampleData.Subcategories()

	sampleData.Generics()

	sampleData.Items()

	deals := sampleData.Deals()

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Approved Transaction Status
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	// Retrieve GUID for Pending Transaction Status
	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Sample transaction, dealCashback
	dealCashback := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	transaction := sampleData.Transaction(user.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 30.00)

	sampleData.DealCashbackTransactionWithCompletedStatus(dealCashback.GUID, user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/cashout_transactions", TestServer.URL, user.GUID)

	cashoutTransactionData := CreateCashoutTransaction{
		Amount:                25.00,
		BankAccountHolderName: "Muhammad Faiz",
		BankAccountNumber:     "1234567890",
		BankCountry:           "Malaysia",
		BankName:              "Maybank",
	}

	jsonBytes, _ := json.Marshal(&cashoutTransactionData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	cashoutTransaction := data["cashout_transaction"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.NotEmpty(t, cashoutTransaction["guid"])
	require.Equal(t, user.GUID, cashoutTransaction["user_guid"])
	require.Equal(t, data["guid"], cashoutTransaction["transaction_guid"])
	require.Equal(t, cashoutTransactionData.BankAccountHolderName, cashoutTransaction["bank_account_name"])
	require.Equal(t, cashoutTransactionData.BankAccountNumber, cashoutTransaction["bank_account_number"])
	require.Equal(t, cashoutTransactionData.BankName, cashoutTransaction["bank_name"])
	require.Equal(t, cashoutTransactionData.BankCountry, cashoutTransaction["bank_country"])
	require.Nil(t, cashoutTransaction["receipt_image"])
	require.Nil(t, cashoutTransaction["transfer_date"])
}
