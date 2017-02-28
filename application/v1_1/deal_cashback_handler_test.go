package v1_1

import (
	"fmt"
	"testing"
	"time"

	"encoding/json"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

func TestCreateDealCashbackShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errors["title"])
}

func TestCreateDealCashbackShouldReturnValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: "",
		DealGUID:         "",
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Validation failed.", errors["title"])
}

func TestCreateDealCashbackShouldReturnShoppingListNotFoundError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: Helper.GenerateUUID(),
		DealGUID:         "001843af-94c1-403d-bb7f-dc7c000defbd",
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 404, status)
	assert.Equal(t, "Shopping List not exists.", errors["title"])
}

func TestCreateDealCashbackShouldReturnDealNotValidDueToDealNotPublish(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	deals := sampleData.Deals()

	// Make deal not valid by setting deal status to draft
	DB.Model(&Ads{}).Where(&Ads{GUID: deals[0].GUID}).Select("status").Update(map[string]interface{}{"status": "draft"})

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Deal already expired or not valid.", errors["title"])
}

func TestCreateDealCashbackShouldReturnDealNotValidDueToExceededQuota(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingListsForUser1 := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	shoppingListsForUser2 := sampleData.ShoppingLists(users[1].GUID, occasions[1].GUID)

	deals := sampleData.Deals()

	dealCashback := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[1].GUID,
		ShoppingListGUID: shoppingListsForUser2[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	DB.Create(&dealCashback)

	// Make deal not valid by setting the deal out of quota
	DB.Model(&Ads{}).Where(&Ads{GUID: deals[0].GUID}).Select("quota").Update(map[string]interface{}{"quota": 1})

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingListsForUser1[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Deal already expired or not valid.", errors["title"])
}

func TestCreateDealCashbackShouldReturnDealNotValidDueToExceededDealLimitPerUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	deals := sampleData.Deals()

	dealCashback := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[1].GUID,
		DealGUID:         deals[0].GUID,
	}

	DB.Create(&dealCashback)

	// Make deal not valid by setting the deal limit per user to 1
	DB.Model(&Ads{}).Where(&Ads{GUID: deals[0].GUID}).Select("perlimit").Update(map[string]interface{}{"perlimit": 1})

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Deal already expired or not valid.", errors["title"])
}

func TestCreateDealCashbackShouldReturnDealNotValidDueToExceededDealExpired(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	deals := sampleData.Deals()

	dealExpiredDate := time.Now().UTC().Add(time.Hour * -24 * -7 * -2)

	// Make the deal expired more than 7 days
	DB.Model(&Ads{}).Select("end_date").Where(&Ads{GUID: deals[0].GUID}).Update(map[string]interface{}{"end_date": dealExpiredDate})

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 422, status)
	assert.Equal(t, "Deal already expired or not valid.", errors["title"])
}

func TestCreateDealCashbackShouldFailedDueToUserAlreadyAddedSameDealIntoShoppingList(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	deals := sampleData.Deals()

	dealCashback := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	DB.Create(&dealCashback)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 409, status)
	assert.Equal(t, "Failed to add deal into the shopping list.", errors["title"])
}

func TestCreateDealCashbackShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	sampleData.Categories()

	sampleData.Subcategories()

	sampleData.Generics()

	sampleData.Items()

	deals := sampleData.Deals()

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	postData := CreateDealCashback{
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	responseData := body.(map[string]interface{})["data"].(map[string]interface{})

	assert.Equal(t, 200, status)
	assert.Equal(t, "Successfully add deal guid "+deals[0].GUID+" to list.", responseData["message"])

	createdShoppingList := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: users[0].GUID, ShoppingListGUID: shoppingLists[0].GUID, Name: deals[0].Name}).
		First(&createdShoppingList)

	cashbackAmount := deals[0].CashbackAmount
	assert.Equal(t, "Drinks", createdShoppingList.Category)
	assert.Equal(t, "Carbonated Drink", createdShoppingList.SubCategory)
	assert.Equal(t, 1, createdShoppingList.AddedFromDeal)
	assert.Equal(t, 1, createdShoppingList.Quantity)
	assert.Equal(t, &cashbackAmount, createdShoppingList.CashbackAmount)
	assert.Equal(t, 0, createdShoppingList.AddedToCart)
	assert.Nil(t, createdShoppingList.DealExpired)

}

func TestViewUserDealCashbacksShouldGiveAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errors["title"])
}

func TestViewUserDealCashbacksShouldGiveAccessTokenBelongsToOtherUserError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, "e5611404-bdce-5943-b5f4-8567b82332c2")

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Your access token belong to other user", errors["title"])
}

func TestViewUserDealCashbacksShouldRemoveDealCashbackWhenExpiredMoreThan7days(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	occasions := sampleData.Occasions()

	users := sampleData.Users()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	deals := sampleData.Deals()

	dealCashback := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	DB.Create(&dealCashback)

	dealCashbackTransactionGUID := Helper.GenerateUUID()

	dealCashback1 := DealCashback{
		GUID:                        Helper.GenerateUUID(),
		UserGUID:                    users[0].GUID,
		ShoppingListGUID:            shoppingLists[1].GUID,
		DealGUID:                    deals[1].GUID,
		DealCashbackTransactionGUID: &dealCashbackTransactionGUID,
	}

	DB.Create(&dealCashback1)

	dealExpiredDate := time.Now().UTC().Add(time.Hour * -24 * -7 * -2)

	// Make the deal expired more than 7 days
	DB.Model(&Ads{}).Select("end_date").Where(&Ads{GUID: deals[0].GUID}).Update(map[string]interface{}{"end_date": dealExpiredDate})

	jwt := systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks?page_number=1&page_limit=1&transaction_status=empty", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	response := body.(map[string]interface{})["data"]

	assert.Equal(t, 200, status)
	assert.Empty(t, response)
}

func TestViewUserDealCashbacksShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	occasions := sampleData.Occasions()

	users := sampleData.Users()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	deals := sampleData.Deals()

	dealCashback1 := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[0].GUID,
		DealGUID:         deals[0].GUID,
	}

	DB.Create(&dealCashback1)

	dealCashbackTransactionGUID := Helper.GenerateUUID()

	dealCashback2 := DealCashback{
		GUID:                        Helper.GenerateUUID(),
		UserGUID:                    users[0].GUID,
		ShoppingListGUID:            shoppingLists[1].GUID,
		DealGUID:                    deals[1].GUID,
		DealCashbackTransactionGUID: &dealCashbackTransactionGUID,
	}

	DB.Create(&dealCashback2)

	dealCashback3 := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[1].GUID,
		DealGUID:         deals[2].GUID,
	}

	DB.Create(&dealCashback3)

	dealCashback4 := DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         users[0].GUID,
		ShoppingListGUID: shoppingLists[2].GUID,
		DealGUID:         deals[3].GUID,
	}

	DB.Create(&dealCashback4)

	DB.Where(&ShoppingList{GUID: shoppingLists[2].GUID}).Delete(&ShoppingList{})

	dealExpiredDate := time.Now().UTC().Add(time.Hour * -24 * -7 * -2)

	// Make the deal expired more than 7 days
	DB.Model(&Ads{}).Select("end_date").Where(&Ads{GUID: deals[0].GUID}).Update(map[string]interface{}{"end_date": dealExpiredDate})

	jwt := systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks?page_number=1&page_limit=1&transaction_status=empty", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwtToken.Token)

	response := body.(map[string]interface{})["data"]

	shoppingListResponse := response.([]interface{})

	pretty.Println(response)

	// Assert Shopping List
	assert.Equal(t, 200, status)
	assert.Len(t, response.([]interface{}), 2)

	assert.Equal(t, "Birthday Party", shoppingListResponse[0].(map[string]interface{})["name"])
	assert.Nil(t, shoppingListResponse[0].(map[string]interface{})["deleted_at"])
	assert.Equal(t, users[0].GUID, shoppingListResponse[0].(map[string]interface{})["user_guid"])

	assert.Equal(t, "Deleted Shopping List", shoppingListResponse[1].(map[string]interface{})["name"])

	dealCashbacksResponse1 := shoppingListResponse[0].(map[string]interface{})["deal_cashbacks"].([]interface{})
	dealCashbacksResponse2 := shoppingListResponse[1].(map[string]interface{})["deal_cashbacks"].([]interface{})

	// Assert Deal Cashbacks
	assert.Len(t, dealCashbacksResponse1, 1)
	assert.Equal(t, users[0].GUID, dealCashbacksResponse1[0].(map[string]interface{})["user_guid"])
	assert.Equal(t, shoppingLists[1].GUID, dealCashbacksResponse1[0].(map[string]interface{})["shopping_list_guid"])
	assert.Equal(t, deals[2].GUID, dealCashbacksResponse1[0].(map[string]interface{})["deal_guid"])
	assert.NotEmpty(t, dealCashbacksResponse1[0].(map[string]interface{})["deal"])
	assert.Equal(t, deals[2].GUID, dealCashbacksResponse1[0].(map[string]interface{})["deal"].(map[string]interface{})["guid"])

	assert.Len(t, dealCashbacksResponse2, 1)
	assert.Equal(t, users[0].GUID, dealCashbacksResponse2[0].(map[string]interface{})["user_guid"])
	assert.Equal(t, shoppingLists[2].GUID, dealCashbacksResponse2[0].(map[string]interface{})["shopping_list_guid"])
	assert.Equal(t, deals[3].GUID, dealCashbacksResponse2[0].(map[string]interface{})["deal_guid"])
	assert.NotEmpty(t, dealCashbacksResponse2[0].(map[string]interface{})["deal"])
	assert.Equal(t, deals[3].GUID, dealCashbacksResponse2[0].(map[string]interface{})["deal"].(map[string]interface{})["guid"])

}
