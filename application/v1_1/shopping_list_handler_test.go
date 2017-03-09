package v1_1

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"
)

func TestViewAllUserShoppingListsShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestViewAllUserShoppingListsWithoutRelationshipShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	DB.Where(&ShoppingList{GUID: shoppingLists[0].GUID}).Delete(&ShoppingList{})

	sampleData.ShoppingLists(users[1].GUID, occasions[1].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	userShoppingLists := body.(map[string]interface{})["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Len(t, userShoppingLists, 2)
	require.Equal(t, users[0].GUID, userShoppingLists[0].(map[string]interface{})["user_guid"])
	require.Equal(t, users[0].GUID, userShoppingLists[1].(map[string]interface{})["user_guid"])
}

func TestViewAllUserShoppingListsWithRelationshipShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingLists := sampleData.ShoppingLists(users[0].GUID, occasions[0].GUID)

	DB.Where(&ShoppingList{GUID: shoppingLists[0].GUID}).Delete(&ShoppingList{})

	sampleData.ShoppingListItem(users[0].GUID, shoppingLists[1].GUID, "Test Shopping List Item", 1)

	sampleData.ShoppingLists(users[1].GUID, occasions[1].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists?include=occasions,items", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	userShoppingLists := body.(map[string]interface{})["data"].([]interface{})

	require.Equal(t, 200, status)
	require.Len(t, userShoppingLists, 2)
	require.Equal(t, users[0].GUID, userShoppingLists[0].(map[string]interface{})["user_guid"])
	require.Equal(t, users[0].GUID, userShoppingLists[1].(map[string]interface{})["user_guid"])

	// Assert relationship data
	require.NotEmpty(t, userShoppingLists[0].(map[string]interface{})["occasions"])
	require.NotEmpty(t, userShoppingLists[0].(map[string]interface{})["items"])
}

func TestCreateNewShoppingListShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, Helper.GenerateUUID())

	postData := CreateShoppingList{
		OccasionGUID: "f1f264aa-c9fd-5801-b92b-e7d6f8449cc9",
		Name:         "Test Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestCreateNewShoppingListShouldReturnValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, users[0].GUID)

	postData := CreateShoppingList{
		OccasionGUID: "",
		Name:         "",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["name"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["occasion_guid"])
}

func TestCreateNewShoppingListShouldReturnOccasionGUIDNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, users[0].GUID)

	postData := CreateShoppingList{
		OccasionGUID: Helper.GenerateUUID(),
		Name:         "Test Shopping Lists",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "Occasion not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestCreateNewShoppingListShouldReturnDuplicateError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, users[0].GUID)

	postData := CreateShoppingList{
		OccasionGUID: occasions[0].GUID,
		Name:         "Test Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 409, status)
	require.Equal(t, "Shopping List already exists.", errors["title"])
}

func TestCreateNewShoppingListShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists", TestServer.URL, users[0].GUID)

	postData := CreateShoppingList{
		OccasionGUID: occasions[0].GUID,
		Name:         "Test Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.NotEmpty(t, data["guid"])
	require.Equal(t, users[0].GUID, data["user_guid"])
	require.Equal(t, postData.Name, data["name"])
	require.Equal(t, occasions[0].GUID, data["occasion_guid"])
}

func TestUpdateShoppingListShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID())

	postData := UpdateShoppingList{
		OccasionGUID: "f1f264aa-c9fd-5801-b92b-e7d6f8449cc9",
		Name:         "Test Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestUpdateShoppingListShouldReturnShoppingListNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, Helper.GenerateUUID())

	postData := UpdateShoppingList{
		OccasionGUID: occasions[0].GUID,
		Name:         "Test Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "Shopping List not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestUpdateShoppingListShouldReturnOccasionNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingList.GUID)

	postData := UpdateShoppingList{
		OccasionGUID: Helper.GenerateUUID(),
		Name:         "Update Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "Occasion not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestUpdateShoppingListUsingSameNameShouldReturnDuplicateError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Update Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingList.GUID)

	postData := UpdateShoppingList{
		Name: "Update Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 409, status)
	require.Equal(t, "Shopping List already exists.", errors["title"])
}

func TestUpdateShoppingListNameShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingList.GUID)

	postData := UpdateShoppingList{
		Name: "Update Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, postData.Name, data["name"])
	require.Equal(t, occasions[0].GUID, data["occasion_guid"])
	require.Equal(t, users[0].GUID, data["user_guid"])
}

func TestUpdateShoppingListOccasionShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingList.GUID)

	postData := UpdateShoppingList{
		OccasionGUID: occasions[1].GUID,
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, shoppingList.Name, data["name"])
	require.Equal(t, postData.OccasionGUID, data["occasion_guid"])
	require.Equal(t, users[0].GUID, data["user_guid"])
}

func TestUpdateShoppingListOccasionAndNameShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	shoppingList := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List")

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingList.GUID)

	postData := UpdateShoppingList{
		OccasionGUID: occasions[1].GUID,
		Name:         "Update Shopping List",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, postData.Name, data["name"])
	require.Equal(t, postData.OccasionGUID, data["occasion_guid"])
	require.Equal(t, users[0].GUID, data["user_guid"])
}

func TestDeleteShoppingListShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, Helper.GenerateUUID(), Helper.GenerateUUID())

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	status, _, body := TestHelper.Request("DELETE", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestDeleteShoppingListShouldReturnShoppingListNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, Helper.GenerateUUID())

	status, _, body := TestHelper.Request("DELETE", []byte{}, requestURL, jwt.Token)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "Shopping List not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["guid"])
}

func TestDeleteShoppingListShouldDeleteShoppingListAndShoppingListItemsAndItemImagesBelongsToTheUserOnly(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	occasions := sampleData.Occasions()

	// Shopping List, Shopping List Item, Shopping List Item Image for User 1
	shoppingListForUser1 := sampleData.ShoppingList(users[0].GUID, occasions[0].GUID, "Test Shopping List User 1")

	shoppingListItem1ForUser1 := sampleData.ShoppingListItem(users[0].GUID, shoppingListForUser1.GUID, "Test Shopping List Item 1 For User 1", 0)

	imagePathForItem1ForUser1 := os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/test/images/shopping_list_item_image.png"

	imageForItem1ForUser1, _ := os.Open(imagePathForItem1ForUser1)

	uploadedImageForItem1ForUser1, _ := TestHelper.UploadToAmazonS3("/item_images/", imageForItem1ForUser1)

	shoppingListItemImageForItem1ForUser1 := sampleData.ShoppingListItemImage(users[0].GUID, shoppingListForUser1.GUID, shoppingListItem1ForUser1.GUID, uploadedImageForItem1ForUser1["path"])

	shoppingListItem2ForUser1 := sampleData.ShoppingListItem(users[0].GUID, shoppingListForUser1.GUID, "Test Shopping List Item 2 For User 1", 1)

	// Shopping List, Shopping List Item, Shopping List Item Image for User 2
	shoppingListForUser2 := sampleData.ShoppingList(users[1].GUID, occasions[1].GUID, "Test Shopping List User 2")

	shoppingListItem1ForUser2 := sampleData.ShoppingListItem(users[1].GUID, shoppingListForUser2.GUID, "Test Shopping List Item 1 For User 2", 0)

	imagePathForItem1ForUser2 := os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/test/images/shopping_list_item_image.png"

	imageForItem1ForUser2, _ := os.Open(imagePathForItem1ForUser2)

	uploadedImageForItem1ForUser2, _ := TestHelper.UploadToAmazonS3("/item_images/", imageForItem1ForUser2)

	shoppingListItemImageForItem1ForUser2 := sampleData.ShoppingListItemImage(users[0].GUID, shoppingListForUser2.GUID, shoppingListItem1ForUser2.GUID, uploadedImageForItem1ForUser2["path"])

	shoppingListItem2ForUser2 := sampleData.ShoppingListItem(users[1].GUID, shoppingListForUser2.GUID, "Test Shopping List Item 2 For User 2", 1)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s/shopping_lists/%s", TestServer.URL, users[0].GUID, shoppingListForUser1.GUID)

	status, _, body := TestHelper.Request("DELETE", []byte{}, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	// Assert to make sure shopping list for user 1 deleted.
	require.Equal(t, 200, status)
	require.NotEmpty(t, data["message"])

	deletedShoppingListForUser1 := &ShoppingList{}

	DB.Unscoped().Model(&ShoppingList{}).Where(&ShoppingList{GUID: shoppingListForUser1.GUID}).Find(deletedShoppingListForUser1)

	require.NotNil(t, deletedShoppingListForUser1.DeletedAt)

	// Assert to make sure shopping list for user 2 not deleted.
	undeletedShoppingListForUser2 := &ShoppingList{}

	DB.Unscoped().Model(&ShoppingList{}).Where(&ShoppingList{GUID: shoppingListForUser2.GUID}).Find(undeletedShoppingListForUser2)

	require.Nil(t, undeletedShoppingListForUser2.DeletedAt)

	// Assert to make sure shopping list items for User 1 deleted.
	deletedShoppingListItem1ForUser1 := &ShoppingListItem{}

	DB.Unscoped().Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem1ForUser1.GUID}).Find(deletedShoppingListItem1ForUser1)

	require.NotNil(t, deletedShoppingListItem1ForUser1.DeletedAt)

	deletedShoppingListItem2ForUser1 := &ShoppingListItem{}

	DB.Unscoped().Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem2ForUser1.GUID}).Find(deletedShoppingListItem2ForUser1)

	require.NotNil(t, deletedShoppingListItem2ForUser1.DeletedAt)

	// Assert to make shopping list item images for Item 1 and for User 1 deleted from database and Amazon S3.
	deletedShoppingListItemImageForItem1ForUser1 := &ShoppingListItemImage{}

	DB.Unscoped().Model(&ShoppingListItemImage{}).Where(&ShoppingListItemImage{GUID: shoppingListItemImageForItem1ForUser1.GUID}).Find(deletedShoppingListItemImageForItem1ForUser1)

	require.NotNil(t, deletedShoppingListItemImageForItem1ForUser1.DeletedAt)

	response, _ := http.Get(deletedShoppingListItemImageForItem1ForUser1.URL)

	require.Equal(t, 403, response.StatusCode)

	// Assert to make sure shopping list items for User 2 not deleted.
	deletedShoppingListItem1ForUser2 := &ShoppingListItem{}

	DB.Unscoped().Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem1ForUser2.GUID}).Find(deletedShoppingListItem1ForUser2)

	require.Nil(t, deletedShoppingListItem1ForUser2.DeletedAt)

	deletedShoppingListItem2ForUser2 := &ShoppingListItem{}

	DB.Unscoped().Model(&ShoppingListItem{}).Where(&ShoppingListItem{GUID: shoppingListItem2ForUser2.GUID}).Find(deletedShoppingListItem2ForUser2)

	require.Nil(t, deletedShoppingListItem2ForUser2.DeletedAt)

	// Assert to make shopping list item images for Item 1 and for User 2 not deleted from database and Amazon S3.
	undeletedShoppingListItemImageForItem1ForUser2 := &ShoppingListItemImage{}

	DB.Unscoped().Model(&ShoppingListItemImage{}).Where(&ShoppingListItemImage{GUID: shoppingListItemImageForItem1ForUser2.GUID}).Find(undeletedShoppingListItemImageForItem1ForUser2)

	require.Nil(t, undeletedShoppingListItemImageForItem1ForUser2.DeletedAt)

	response, _ = http.Get(undeletedShoppingListItemImageForItem1ForUser2.URL)

	require.Equal(t, 200, response.StatusCode)
}
