package v1_1

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDealCashbackTransactionShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("POST", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errors["title"])
}

func TestCreateDealCashbackTransactionShouldReturnErrorDealCashbackGUIDsRequired(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": "",
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "", "", jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Validation failed.", errors["title"])
	assert.NotEmpty(t, errors["detail"].(map[string]interface{})["deal_cashback_guids"])
}

func TestCreateDealCashbackTransactionShouldReturnReceiptImageRequired(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": Helper.GenerateUUID() + "," + Helper.GenerateUUID(),
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "", "", jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Validation failed.", errors["title"])
	assert.Equal(t, "The receipt_image parameter is required.", errors["detail"])
}

func TestCreateDealCashbackTransactionShouldReturnErrorInvalidFileType(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": Helper.GenerateUUID() + "," + Helper.GenerateUUID(),
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/files/test_pdf_file.pdf", jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 400, status)
	assert.Equal(t, "Invalid file type.", errors["title"])
}

func TestCreateDealCashbackTransactionShouldReturnErrorDealCashbackGUIDsNotFound(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user.GUID)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": Helper.GenerateUUID() + "," + Helper.GenerateUUID(),
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/receipt_smaller.jpg", jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 404, status)
	assert.Equal(t, "Deal Cashback GUID not exists.", errors["title"])
}

func TestCreateDealCashbackTransactionShouldCreateTransactionRecord(t *testing.T) {
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

	dealCashback1 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	dealCashback2 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[1].GUID, nil)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": dealCashback1.GUID + "," + dealCashback2.GUID,
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/receipt_smaller.jpg", jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	pendingTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "pending").Find(&pendingTransactionStatus)

	dealCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "deal_redemption").Find(&dealCashbackTransactionType)

	assert.Equal(t, 200, status)
	assert.Equal(t, user.GUID, data["user_guid"])
	assert.Equal(t, 0, int(data["read_status"].(interface{}).(float64)))
	assert.Equal(t, deals[0].CashbackAmount+deals[1].CashbackAmount, data["total_amount"])
	assert.NotEmpty(t, data["reference_id"])
	assert.Nil(t, data["deleted_at"])
	assert.Equal(t, pendingTransactionStatus.GUID, data["transaction_status_guid"])
	assert.Equal(t, dealCashbackTransactionType.GUID, data["transaction_type_guid"])
}

func TestCreateDealCashbackTransactionShouldCreateDealCashbackTransactionRecord(t *testing.T) {
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

	dealCashback1 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	dealCashback2 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[1].GUID, nil)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": dealCashback1.GUID + "," + dealCashback2.GUID,
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/receipt_smaller.jpg", jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	dealCashbackTransaction := data["deal_cashback_transaction"].(map[string]interface{})

	assert.Equal(t, 200, status)
	assert.Equal(t, dealCashbackTransaction["user_guid"], data["user_guid"])
	assert.Equal(t, data["guid"], dealCashbackTransaction["transaction_guid"])
	assert.Nil(t, dealCashbackTransaction["deleted_at"])
	assert.Nil(t, dealCashbackTransaction["verification_date"])
	assert.NotEmpty(t, dealCashbackTransaction["guid"])
	assert.Len(t, dealCashbackTransaction["deal_cashbacks"].([]interface{}), 2)

	response, _ := http.Get(dealCashbackTransaction["receipt_url"].(string))
	assert.Equal(t, 200, response.StatusCode)

}

func TestCreateDealCashbackTransactionShouldUpdateTransactionGUIDInDealCashbacks(t *testing.T) {
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

	dealCashback1 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[0].GUID, nil)

	dealCashback2 := sampleData.DealCashback(user.GUID, shoppingList.GUID, deals[1].GUID, nil)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": dealCashback1.GUID + "," + dealCashback2.GUID,
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/receipt_smaller.jpg", jwtToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	dealCashbackTransaction := data["deal_cashback_transaction"].(map[string]interface{})

	updatedDealCashback1 := &DealCashback{}

	DB.Model(&DealCashback{}).Where(&DealCashback{GUID: dealCashback1.GUID}).First(updatedDealCashback1)

	updatedDealCashback2 := &DealCashback{}

	DB.Model(&DealCashback{}).Where(&DealCashback{GUID: dealCashback2.GUID}).First(updatedDealCashback2)

	assert.Equal(t, 200, status)
	assert.Equal(t, dealCashbackTransaction["guid"], *updatedDealCashback1.DealCashbackTransactionGUID)
	assert.Equal(t, dealCashbackTransaction["guid"], *updatedDealCashback2.DealCashbackTransactionGUID)
}

func TestCreateDealCashbackTransactionShouldAddShoppingListItemToCart(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("61111111111", "muhdfaiz@mediacliq.my")

	device := sampleData.DeviceWithUserGUID(user.GUID)

	occasions := sampleData.Occasions()

	shoppingList1 := sampleData.ShoppingList(user.GUID, occasions[0].GUID, "Test Shopping List 1")

	shoppingList2 := sampleData.ShoppingList(user.GUID, occasions[0].GUID, "Test Shopping List 2")

	sampleData.Categories()

	sampleData.Subcategories()

	sampleData.Generics()

	sampleData.Items()

	deals := sampleData.Deals()

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	dealCashback1 := sampleData.DealCashback(user.GUID, shoppingList1.GUID, deals[0].GUID, nil)

	dealCashback2 := sampleData.DealCashback(user.GUID, shoppingList2.GUID, deals[1].GUID, nil)

	shoppingListItemForDealCashback1 := sampleData.ShoppingListItemForDeal(user.GUID, shoppingList1.GUID, deals[0].Name, "Drinks", "Carbonated Drinks", deals[0].GUID, deals[0].CashbackAmount, 0)

	shoppingListItemForDealCashback2 := sampleData.ShoppingListItemForDeal(user.GUID, shoppingList2.GUID, deals[1].Name, "Drinks", "Carbonated Drinks", deals[1].GUID, deals[1].CashbackAmount, 0)

	shoppingListItem1 := sampleData.ShoppingListItem(user.GUID, shoppingList1.GUID, "Test Shopping List Item 1", 0)

	shoppingListItem2 := sampleData.ShoppingListItem(user.GUID, shoppingList2.GUID, "Test Shopping List Item 2", 0)

	jwtToken, _ := JWT.GenerateToken(user.GUID, user.PhoneNo, device.UUID, "")

	dealCashbackTransactionData := map[string]string{
		"deal_cashback_guids": dealCashback1.GUID + "," + dealCashback2.GUID,
	}

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/transactions/deal_cashback_transactions", TestServer.URL, user.GUID)

	status, _, _ := TestHelper.MultipartRequest(requestURL, "POST", dealCashbackTransactionData, "receipt_image",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/receipt_smaller.jpg", jwtToken.Token)

	updatedShoppingListItemForDealCashback1 := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItemForDealCashback1.GUID}).First(updatedShoppingListItemForDealCashback1)

	updatedShoppingListItemForDealCashback2 := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItemForDealCashback2.GUID}).First(updatedShoppingListItemForDealCashback2)

	assert.Equal(t, 200, status)

	// Verify API update added to cart value to 1 or not for shopping list item deal cashback 1 after user
	// successfuly redeem deal cashback.
	assert.Equal(t, 1, updatedShoppingListItemForDealCashback1.AddedToCart)
	assert.Equal(t, 1, updatedShoppingListItemForDealCashback1.AddedFromDeal)
	assert.Equal(t, deals[0].GUID, *updatedShoppingListItemForDealCashback1.DealGUID)
	assert.Nil(t, updatedShoppingListItemForDealCashback1.DealExpired)

	// Verify API update added to cart value to 1 or not for shopping list item deal cashback 2 after user
	// successfuly redeem deal cashback.
	assert.Equal(t, 1, updatedShoppingListItemForDealCashback2.AddedToCart)
	assert.Equal(t, 1, updatedShoppingListItemForDealCashback2.AddedFromDeal)
	assert.Equal(t, deals[1].GUID, *updatedShoppingListItemForDealCashback2.DealGUID)
	assert.Nil(t, updatedShoppingListItemForDealCashback1.DealExpired)

	notUpdateShoppingList1 := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem1.GUID}).First(notUpdateShoppingList1)

	notUpdateShoppingList2 := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem2.GUID}).First(notUpdateShoppingList2)

	// Make Sure API didn't update other shopping list item added to cart value
	assert.Equal(t, 0, notUpdateShoppingList1.AddedToCart)
	assert.Equal(t, 0, notUpdateShoppingList2.AddedToCart)
}
