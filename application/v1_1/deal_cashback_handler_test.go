package v1_1

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"encoding/json"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

func TestCreateDealCashbackShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Access token error", errors["title"])
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

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Validation failed.", errors["title"])
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

	require.Equal(testingT{t}, 404, status)
	require.Equal(testingT{t}, "Shopping List not exists.", errors["title"])
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

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Deal already expired or not valid.", errors["title"])
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

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Deal already expired or not valid.", errors["title"])
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

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Deal already expired or not valid.", errors["title"])
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

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Deal already expired or not valid.", errors["title"])
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

	require.Equal(testingT{t}, 409, status)
	require.Equal(testingT{t}, "Failed to add deal into the shopping list.", errors["title"])
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

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, "Successfully add deal guid "+deals[0].GUID+" to list.", responseData["message"])

	createdShoppingList := &ShoppingListItem{}

	DB.Model(&ShoppingListItem{}).Where(&ShoppingListItem{UserGUID: users[0].GUID, ShoppingListGUID: shoppingLists[0].GUID, Name: deals[0].Name}).
		First(&createdShoppingList)

	cashbackAmount := deals[0].CashbackAmount
	require.Equal(testingT{t}, "Drinks", createdShoppingList.Category)
	require.Equal(testingT{t}, "Carbonated Drink", createdShoppingList.SubCategory)
	require.Equal(testingT{t}, 1, createdShoppingList.AddedFromDeal)
	require.Equal(testingT{t}, 1, createdShoppingList.Quantity)
	require.Equal(testingT{t}, &cashbackAmount, createdShoppingList.CashbackAmount)
	require.Equal(testingT{t}, 0, createdShoppingList.AddedToCart)
	require.Nil(testingT{t}, createdShoppingList.DealExpired)

}

func TestViewUserDealCashbacksShouldGiveAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/deal_cashbacks", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Access token error", errors["title"])
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

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Your access token belong to other user", errors["title"])
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

	require.Equal(testingT{t}, 200, status)
	require.Empty(testingT{t}, response)
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

	// Assert Shopping List
	require.Equal(testingT{t}, 200, status)
	require.Len(testingT{t}, response.([]interface{}), 2)

	require.Equal(testingT{t}, "Birthday Party", shoppingListResponse[0].(map[string]interface{})["name"])
	require.Nil(testingT{t}, shoppingListResponse[0].(map[string]interface{})["deleted_at"])
	require.Equal(testingT{t}, users[0].GUID, shoppingListResponse[0].(map[string]interface{})["user_guid"])

	require.Equal(testingT{t}, "Deleted Shopping List", shoppingListResponse[1].(map[string]interface{})["name"])

	dealCashbacksResponse1 := shoppingListResponse[0].(map[string]interface{})["deal_cashbacks"].([]interface{})
	dealCashbacksResponse2 := shoppingListResponse[1].(map[string]interface{})["deal_cashbacks"].([]interface{})

	// Assert Deal Cashbacks
	require.Len(testingT{t}, dealCashbacksResponse1, 1)
	require.Equal(testingT{t}, users[0].GUID, dealCashbacksResponse1[0].(map[string]interface{})["user_guid"])
	require.Equal(testingT{t}, shoppingLists[1].GUID, dealCashbacksResponse1[0].(map[string]interface{})["shopping_list_guid"])
	require.Equal(testingT{t}, deals[2].GUID, dealCashbacksResponse1[0].(map[string]interface{})["deal_guid"])
	require.NotEmpty(testingT{t}, dealCashbacksResponse1[0].(map[string]interface{})["deal"])
	require.Equal(testingT{t}, deals[2].GUID, dealCashbacksResponse1[0].(map[string]interface{})["deal"].(map[string]interface{})["guid"])

	require.Len(testingT{t}, dealCashbacksResponse2, 1)
	require.Equal(testingT{t}, users[0].GUID, dealCashbacksResponse2[0].(map[string]interface{})["user_guid"])
	require.Equal(testingT{t}, shoppingLists[2].GUID, dealCashbacksResponse2[0].(map[string]interface{})["shopping_list_guid"])
	require.Equal(testingT{t}, deals[3].GUID, dealCashbacksResponse2[0].(map[string]interface{})["deal_guid"])
	require.NotEmpty(testingT{t}, dealCashbacksResponse2[0].(map[string]interface{})["deal"])
	require.Equal(testingT{t}, deals[3].GUID, dealCashbacksResponse2[0].(map[string]interface{})["deal"].(map[string]interface{})["guid"])

}
