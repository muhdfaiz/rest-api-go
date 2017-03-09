package v1_1

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestViewDealCashbackTransactionShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/deal_cashback_transactions", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestViewDealCashbackTransactionShouldUpdateReadStatus(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

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

	// Retrieve GUID for Deal Cashback Transaction Type
	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	// Create Sample transaction, dealCashback
	dealCashback := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	transaction := sampleData.Transaction(user.GUID, dealRedemptionTransactionType.GUID, approvedTransactionStatus.GUID, 0, 30.00)

	sampleData.DealCashbackTransactionWithCompletedStatus(dealCashback.GUID, user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/deal_cashback_transactions", TestServer.URL, user.GUID, transaction.GUID)

	status, _, _ := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	require.Equal(t, 200, status)

	updatedDealCashbackTransaction := &Transaction{}

	DB.Model(&Transaction{}).Where(&Transaction{GUID: transaction.GUID}).First(&updatedDealCashbackTransaction)

	require.Equal(t, 1, updatedDealCashbackTransaction.ReadStatus)
}

func TestViewDealCashbackTransactionShouldReturnTransactionDetailIncludingRelationship(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

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

	// Retrieve GUID for Deal Cashback Transaction Type
	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	// Create Sample transaction, dealCashback
	dealCashback1 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	dealCashback2 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[1].GUID, nil)

	transaction := sampleData.Transaction(user.GUID, dealRedemptionTransactionType.GUID, approvedTransactionStatus.GUID, 0, deals[0].CashbackAmount+deals[1].CashbackAmount)

	dealCashbackTransaction1 := sampleData.DealCashbackTransactionWithCompletedStatus(dealCashback1.GUID, user.GUID, transaction.GUID)

	DB.Model(&DealCashback{}).Where("guid IN (?)", []string{dealCashback1.GUID, dealCashback2.GUID}).Select("deal_cashback_transaction_guid").Updates(map[string]interface{}{
		"deal_cashback_transaction_guid": dealCashbackTransaction1.GUID,
	})

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/deal_cashback_transactions", TestServer.URL, user.GUID, transaction.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	dealCashbackGroupByShoppingList := data["deal_cashback_transaction"].(map[string]interface{})["deal_cashbacks_group_by_shopping_list"].([]interface{})

	dealCashbacks := dealCashbackGroupByShoppingList[0].(map[string]interface{})["deal_cashbacks"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, transaction.GUID, data["guid"])
	require.NotEmpty(t, data["deal_cashback_transaction"])
	require.Equal(t, 2, int(data["deal_cashback_transaction"].(map[string]interface{})["total_deal"].(float64)))
	require.NotEmpty(t, data["transaction_status"])
	require.NotEmpty(t, data["transaction_type"])
	require.NotEmpty(t, dealCashbackGroupByShoppingList)
	require.Len(t, dealCashbackGroupByShoppingList, 1)
	require.Equal(t, shoppingList.GUID, dealCashbackGroupByShoppingList[0].(map[string]interface{})["guid"])
	require.Equal(t, shoppingList.Name, dealCashbackGroupByShoppingList[0].(map[string]interface{})["name"])
	require.Equal(t, user.GUID, dealCashbackGroupByShoppingList[0].(map[string]interface{})["user_guid"])
	require.NotEmpty(t, dealCashbackGroupByShoppingList[0].(map[string]interface{})["deal_cashbacks"])
	require.Len(t, dealCashbacks, 2)
	require.NotEmpty(t, dealCashbacks[0].(map[string]interface{})["deal"], 2)
}

func TestViewCashoutTransactionShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/cashout_transactions", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestViewCashoutTransactionShouldNotUpdateReadStatusWhenTransactionStatusPending(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Pending Transaction Status
	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	transaction := sampleData.Transaction(user.GUID, cashoutTransactionType.GUID, pendingTransactionStatus.GUID, 0, 30.00)

	sampleData.CashoutTransaction(user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/cashout_transactions", TestServer.URL, user.GUID, transaction.GUID)

	status, _, _ := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	require.Equal(t, 200, status)

	updatedCashoutTransaction := &Transaction{}

	DB.Model(&Transaction{}).Where(&Transaction{GUID: transaction.GUID}).First(&updatedCashoutTransaction)

	require.Equal(t, 0, updatedCashoutTransaction.ReadStatus)
}

func TestViewCashoutTransactionShouldUpdateReadStatusWhenTransactionStatusNotPending(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Approved Transaction Status
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&approvedTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	transaction := sampleData.Transaction(user.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 30.00)

	sampleData.CashoutTransaction(user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/cashout_transactions", TestServer.URL, user.GUID, transaction.GUID)

	status, _, _ := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	require.Equal(t, 200, status)

	updatedCashoutTransaction := &Transaction{}

	DB.Model(&Transaction{}).Where(&Transaction{GUID: transaction.GUID}).First(&updatedCashoutTransaction)

	require.Equal(t, 1, updatedCashoutTransaction.ReadStatus)
}

func TestViewCashoutTransactionShouldReturnTransactionDetailIncludingRelationship(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 30.00)

	device := sampleData.DeviceWithUserGUID(user.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	// Retrieve GUID for Pending Transaction Status
	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	// Retrieve GUID for Cashout Transaction Type
	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	transaction := sampleData.Transaction(user.GUID, cashoutTransactionType.GUID, pendingTransactionStatus.GUID, 0, 30.00)

	sampleData.CashoutTransaction(user.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/cashout_transactions", TestServer.URL, user.GUID, transaction.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	cashoutTransactionData := data["cashout_transaction"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, transaction.GUID, data["guid"])
	require.Equal(t, cashoutTransactionType.GUID, data["transaction_type_guid"])
	require.Equal(t, pendingTransactionStatus.GUID, data["transaction_status_guid"])
	require.Equal(t, user.GUID, data["user_guid"])
	require.Equal(t, user.GUID, cashoutTransactionData["user_guid"])
	require.Equal(t, transaction.GUID, cashoutTransactionData["transaction_guid"])
	require.NotEmpty(t, data["transaction_status"])
	require.NotEmpty(t, data["transaction_type"])
}

func TestViewReferralCashbackTransactionShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/referral_cashback_transactions", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestViewReferralCashbackTransactionShouldReturnShouldReturnTransactionDetailIncludingRelationship(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 30.00)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID for Approved Transaction Status
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	// Retrieve GUID for Referral Cashback Transaction Type
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	referralPriceSetting := &Setting{}

	DB.Model(&Setting{}).Where("slug = ?", "referral_price").First(&referralPriceSetting)

	referralPriceValue, _ := strconv.ParseFloat(referralPriceSetting.Value, 64)

	transaction := sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, approvedTransactionStatus.GUID, 0, referralPriceValue)

	sampleData.ReferralCashbackTransaction(user1.GUID, user2.GUID, transaction.GUID)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/%s/referral_cashback_transactions", TestServer.URL, user1.GUID, transaction.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	referralCashbackTransactionData := data["referral_cashback_transaction"].(map[string]interface{})

	referrerData := referralCashbackTransactionData["referrer"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, transaction.GUID, data["guid"])
	require.Equal(t, referralCashbackTransactionType.GUID, data["transaction_type_guid"])
	require.Equal(t, approvedTransactionStatus.GUID, data["transaction_status_guid"])
	require.Equal(t, user1.GUID, data["user_guid"])
	require.Equal(t, user1.GUID, referralCashbackTransactionData["user_guid"])
	require.Equal(t, transaction.GUID, referralCashbackTransactionData["transaction_guid"])
	require.Equal(t, user2.GUID, referrerData["guid"])
	require.NotEmpty(t, data["transaction_status"])
	require.NotEmpty(t, data["transaction_type"])
	require.Equal(t, user2.GUID, referrerData["guid"])
}

func TestViewUserTransactionsShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?is_read=1&transaction_status=pending,approved&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=5", TestServer.URL, Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestViewUserTransactionShouldReturnAllUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=5", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 4, int(data["total_data"].(float64)))
	require.Len(t, innerData, 4)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[1].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[2].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[3].(map[string]interface{})["user_guid"])
	require.NotEmpty(t, innerData[0].(map[string]interface{})["transaction_type"])
	require.NotEmpty(t, innerData[0].(map[string]interface{})["transaction_status"])
}

func TestPageNumberInViewUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?include=Transactiontypes,Transactionstatuses&page_number=4&page_limit=1", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 4, int(data["total_data"].(float64)))
	require.Len(t, innerData, 2)
}

func TestPageLimitInViewUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=3", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 4, int(data["total_data"].(float64)))
	require.Len(t, innerData, 3)
}

func TestViewUnreadUserTransactionShouldReturnUnreadUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?is_read=0&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 2, int(data["total_data"].(float64)))
	require.Len(t, innerData, 2)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[1].(map[string]interface{})["user_guid"])
	require.Equal(t, 0, int(innerData[0].(map[string]interface{})["read_status"].(float64)))
	require.Equal(t, 0, int(innerData[1].(map[string]interface{})["read_status"].(float64)))
}

func TestViewReadUserTransactionShouldReturnReadUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?is_read=1&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 2, int(data["total_data"].(float64)))
	require.Len(t, innerData, 2)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[1].(map[string]interface{})["user_guid"])
	require.Equal(t, 1, int(innerData[0].(map[string]interface{})["read_status"].(float64)))
	require.Equal(t, 1, int(innerData[1].(map[string]interface{})["read_status"].(float64)))
}

func TestViewPendingUserTransactionShouldReturnPendingUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?transaction_status=pending&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 1, int(data["total_data"].(float64)))
	require.Len(t, innerData, 1)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, pendingTransactionStatus.Slug, innerData[0].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])
}

func TestViewApprovedUserTransactionShouldReturnApprovedUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?transaction_status=approved&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 1, int(data["total_data"].(float64)))
	require.Len(t, innerData, 1)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, approvedTransactionStatus.Slug, innerData[0].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])
}

func TestViewRejectUserTransactionShouldReturnRejectUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?transaction_status=reject&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 1, int(data["total_data"].(float64)))
	require.Len(t, innerData, 1)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, rejectTransactionStatus.Slug, innerData[0].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])
}

func TestViewPartialSuccessUserTransactionShouldReturnPartialSuccessUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?transaction_status=partial_success&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 1, int(data["total_data"].(float64)))
	require.Len(t, innerData, 1)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, partialSuccessTransactionStatus.Slug, innerData[0].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])
}

func TestViewPendingAndUnreadUserTransactionShouldReturnEmptyPendingAndUnreadUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?is_read=0&transaction_status=pending&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Empty(t, innerData)
}

func TestViewPendingAndApprovedAndReadUserTransactionShouldReturnPendingAndApprovedAndReadUserTransaction(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user1 := sampleData.UserWithCustomWalletAmount("61111111111", "muhdfaiz@mediacliq.my", 25.10)

	user2 := sampleData.User("622222222222", "muhdfaiz2@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user1.GUID)

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	sampleData.Settings("true", "5", "3")

	// Retrieve GUID Transaction Statuses
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	partialSuccessTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "partial_success").Find(&partialSuccessTransactionStatus)

	rejectTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "reject").Find(&rejectTransactionStatus)

	// Retrieve GUID for Transaction Types
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	dealRedemptionTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealRedemptionTransactionType)

	cashoutTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "cashout").Find(&cashoutTransactionType)

	// Create Pending Transaction for User 1
	sampleData.Transaction(user1.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 1, 5.00)

	// Create Approved Transaction for User 1
	sampleData.Transaction(user1.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 1, 27.50)

	// Create Reject Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 0, 1.70)

	// Create Partial Success Transaction for User 1
	sampleData.Transaction(user1.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 0, 3.90)

	// Create Pending Transaction for User 2
	sampleData.Transaction(user2.GUID, referralCashbackTransactionType.GUID, pendingTransactionStatus.GUID, 0, 1.70)

	// Create Approved Transaction for User 2
	sampleData.Transaction(user2.GUID, cashoutTransactionType.GUID, approvedTransactionStatus.GUID, 0, 29.40)

	// Create Reject Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, rejectTransactionStatus.GUID, 1, 2.30)

	// Create Partial Success Transaction for User 2
	sampleData.Transaction(user2.GUID, dealRedemptionTransactionType.GUID, partialSuccessTransactionStatus.GUID, 1, 3.40)

	jwtToken, _ := JWT.GenerateToken(user1.GUID, user1.PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions?is_read=1&transaction_status=pending,approved&include=Transactiontypes,Transactionstatuses&page_number=1&page_limit=2", TestServer.URL, user1.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	innerData := data["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, 2, int(data["total_data"].(float64)))
	require.Len(t, innerData, 2)
	require.Equal(t, user1.GUID, innerData[0].(map[string]interface{})["user_guid"])
	require.Equal(t, user1.GUID, innerData[1].(map[string]interface{})["user_guid"])
	require.Equal(t, pendingTransactionStatus.Slug, innerData[0].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])
	require.Equal(t, approvedTransactionStatus.Slug, innerData[1].(map[string]interface{})["transaction_status"].(map[string]interface{})["slug"])

}
