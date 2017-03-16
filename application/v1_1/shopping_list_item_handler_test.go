package v1_1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewShoppingListItemShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s/items/%s", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID(), Helper.GenerateUUID())

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errors["title"])
}

func TestViewShoppingListItemShouldReturnShoppingListNotExist(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s/items/%s", TestServer.URL, users[0].GUID, Helper.GenerateUUID(), Helper.GenerateUUID())

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 404, status)
	assert.Equal(t, "Shopping List not exists.", errors["title"])
	assert.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestViewShoppingListItemShouldReturnShoppingListItemNotExist(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping Lists")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s/items/%s", TestServer.URL, users[0].GUID, shoppingList.GUID, Helper.GenerateUUID())

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 404, status)
	assert.Equal(t, "Shopping List Item not exists.", errors["title"])
	assert.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestViewShoppingListItemWithoutRelationShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping Lists")

	shoppingListItem := sampleData.ShoppingListItem(users[0].GUID, shoppingList.GUID, "Test Shopping List Item", 0)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s/items/%s", TestServer.URL, users[0].GUID, shoppingList.GUID, shoppingListItem.GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	assert.Equal(t, 200, status)
	assert.Equal(t, shoppingListItem.GUID, data["guid"])
	assert.Equal(t, users[0].GUID, data["user_guid"])
	assert.Equal(t, shoppingList.GUID, data["shopping_list_guid"])
	assert.Equal(t, shoppingListItem.Name, data["name"])
	assert.Nil(t, data["cashback_amount"])
	assert.Equal(t, 0.00, data["added_from_deal"])
}
