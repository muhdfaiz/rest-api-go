package v1_1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TestCreateUserShouldReturnNumericValidationErrors function used to test if API will return error or not if the input
// for numeric field like facebook_id not numeric.
func TestCreateUserShouldReturnNumericValidationErrors(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"facebook_id":       "a",
		"phone_no":          "a",
		"verification_code": "a",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})
	errorDetail := errorData["detail"].(map[string]interface{})

	require.Equal(testingT{t}, 422, status, "Status code must be 422")
	require.NotEmpty(testingT{t}, errorDetail["facebook_id"])
	require.NotEmpty(testingT{t}, errorDetail["phone_no"])
	require.NotEmpty(testingT{t}, errorDetail["verification_code"])
}

// TestCreateUserShouldReturnRequiredFieldValidationErrors function used to test if API will return error if the input data not
// contain all the required fields.
func TestCreateUserShouldReturnRequiredFieldValidationErrors(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"facebook_id": "123123123123",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})
	errorDetail := errorData["detail"].(map[string]interface{})

	require.Equal(testingT{t}, 422, status, "Status code must be 422")
	require.NotEmpty(testingT{t}, errorDetail["name"])
	require.NotEmpty(testingT{t}, errorDetail["email"])
	require.NotEmpty(testingT{t}, errorDetail["phone_no"])
	require.NotEmpty(testingT{t}, errorDetail["verification_code"])
	require.NotEmpty(testingT{t}, errorDetail["device_uuid"])
}

// TestPhoneNumberShouldUniqueDuringCreateUser function used to test if API return error or not when the input data
// contain the same phone number that already exist in database.
func TestPhoneNumberShouldUniqueDuringCreateUser(t *testing.T) {
	sampleData := SampleData{DB: DB}

	users := sampleData.Users()
	fmt.Println(users)
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60121234567",
		"device_uuid":       "F4EE037F5340586BE625EE7678797DD0",
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 409, status, "Status code must be 409")
	require.Equal(testingT{t}, "Duplicate entry 'phone_no' for key '"+userData["phone_no"]+"'.", errorData["detail"].(map[string]interface{})["message"])
}

// TestFacebookIDMustValid function used to test if API will validate the facebook id
// valid or not.
func TestFacebookIDMustValidDuringCreateUser(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"facebook_id":       "111111111111",
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60131234567",
		"device_uuid":       "F4EE037F5340586BE625EE7678797DD0",
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 400, status, "Status code must be 409")
	require.NotEmpty(testingT{t}, errorData["detail"].(map[string]interface{})["facebook_id"])
}

// TestFacebookIDMustUnique function used to check if API will return error or not
// if the input data contain facebook id that already exist in database.
func TestFacebookIDMustUniqueDuringCreateUser(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"facebook_id":       "100013413336774",
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60131234567",
		"device_uuid":       "F4EE037F5340586BE625EE7678797DD0",
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 409, status, "Status code must be 409")
	require.Equal(testingT{t}, "Duplicate entry 'facebook_id' for key '100013413336774'.", errorData["detail"].(map[string]interface{})["message"])
}

func TestReferralCodeMustValidDuringCreateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Settings("true", "2", "5")

	device, _ := sampleData.DeviceWithoutUserGUID()

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"referral_code":     "AAA12345",
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60131234567",
		"device_uuid":       device.UUID,
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 404, status, "Status code must be 404")
	require.NotEmpty(testingT{t}, errorData["detail"].(map[string]interface{})["referral_code"])
}

func TestCreateUserWithReferralCodeNotActiveAndDeviceUUIDMustValid(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Settings("false", "2", "5")

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"referral_code":     "AAA12345",
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60131234567",
		"device_uuid":       "123asdasd213213",
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 404, status, "Status code must be 404")
	require.NotEmpty(testingT{t}, errorData["detail"].(map[string]interface{})["uuid"])
}

// TestDebug function used to test if API can create new user and generate access token without valid verification code.
func TestCreateUserWithDebugMode(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	requestURL := fmt.Sprintf("%s/v1_1/users?debug=1", TestServer.URL)

	userData := map[string]string{
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60174862127",
		"device_uuid":       device.UUID,
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.NotEmpty(testingT{t}, data["access_token"])
	require.NotEmpty(testingT{t}, data["user"])
}

// TestDebugToken function used to test if API will generate access token based on debug token value.
func TestCreateUserWithDebugToken(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	requestURL := fmt.Sprintf("%s/v1_1/users?debug=1&debug_token=6", TestServer.URL)

	userData := map[string]string{
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60174862127",
		"device_uuid":       device.UUID,
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	minutes := time.Minute * time.Duration(6)

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, time.Now().UTC().Add(minutes).Format(time.RFC3339), data["access_token"].(map[string]interface{})["expired"])
	require.NotEmpty(testingT{t}, data["user"])
}

func TestCreateUserWithReferral(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Settings("true", "2", "5")

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	device, _ := sampleData.DeviceWithoutUserGUID()

	users := sampleData.Users()

	referralPriceSetting := &Setting{}

	// Retrieve GUID for Approved Transaction Status
	approvedTransactionStatus := &TransactionStatus{}

	DB.Model(&TransactionStatus{}).Where("slug = ?", "approved").Find(&approvedTransactionStatus)

	// Retrieve GUID for Referral Cashback Transaction Type
	referralCashbackTransactionType := &TransactionType{}

	DB.Model(&TransactionType{}).Where("slug = ?", "referral_cashback").Find(&referralCashbackTransactionType)

	DB.Model(&Setting{}).Where("slug = ?", "referral_price").Find(&referralPriceSetting)

	requestURL := fmt.Sprintf("%s/v1_1/users?debug=1", TestServer.URL)

	userData := map[string]string{
		"referral_code":     users[0].ReferralCode,
		"name":              "Muhammad Faiz",
		"email":             "muhdfaiz@mediacliq.my",
		"phone_no":          "60174862127",
		"device_uuid":       device.UUID,
		"verification_code": "4544",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	newUser := data["user"].(map[string]interface{})

	user1 := &User{}

	DB.Model(&User{}).Where("guid = ?", users[0].GUID).Find(&user1)

	require.Equal(testingT{t}, referralPriceSetting.Value, strconv.FormatFloat(user1.Wallet, 'f', 0, 64))

	referralCashbackTransaction := &ReferralCashbackTransaction{}

	DB.Model(&ReferralCashbackTransaction{}).Where("user_guid = ?", user1.GUID).Find(&referralCashbackTransaction)

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, user1.GUID, referralCashbackTransaction.UserGUID)
	require.Equal(testingT{t}, newUser["guid"], referralCashbackTransaction.ReferrerGUID)
	require.NotEmpty(testingT{t}, referralCashbackTransaction.TransactionGUID)

	transaction := &Transaction{}

	DB.Model(&Transaction{}).Where("guid = ?", referralCashbackTransaction.TransactionGUID).Find(&transaction)

	referralPriceInFloat, _ := strconv.ParseFloat(referralPriceSetting.Value, 64)

	require.NotEmpty(testingT{t}, transaction.GUID)
	require.Equal(testingT{t}, user1.GUID, transaction.UserGUID)
	require.Equal(testingT{t}, referralCashbackTransactionType.GUID, transaction.TransactionTypeGUID)
	require.Equal(testingT{t}, approvedTransactionStatus.GUID, transaction.TransactionStatusGUID)
	require.NotEmpty(testingT{t}, transaction.ReferenceID)
	require.Equal(testingT{t}, referralPriceInFloat, transaction.TotalAmount)

}

func TestMaxReferralPerUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Settings("true", "2", "2")

	sampleData.TransactionStatuses()

	sampleData.TransactionTypes()

	devices, _ := sampleData.DevicesWithoutUserGUID()

	referralPriceSetting := &Setting{}

	DB.Model(&Setting{}).Where("slug = ?", "referral_price").Find(&referralPriceSetting)

	users := sampleData.Users()

	requestURL := fmt.Sprintf("%s/v1_1/users?debug=1", TestServer.URL)

	userData1 := map[string]string{
		"referral_code":     users[0].ReferralCode,
		"name":              "Test User 1",
		"email":             "testuser1@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       devices[0].UUID,
		"verification_code": "1111",
	}

	jsonBytes1, _ := json.Marshal(userData1)

	TestHelper.Request("POST", jsonBytes1, requestURL, "")

	userData2 := map[string]string{
		"referral_code":     users[0].ReferralCode,
		"name":              "Test User 2",
		"email":             "testuser2@mediacliq.my",
		"phone_no":          "60122222222",
		"device_uuid":       devices[1].UUID,
		"verification_code": "2222",
	}

	jsonBytes2, _ := json.Marshal(userData2)

	TestHelper.Request("POST", jsonBytes2, requestURL, "")

	userData3 := map[string]string{
		"referral_code":     users[0].ReferralCode,
		"name":              "Test User 3",
		"email":             "testuser3@mediacliq.my",
		"phone_no":          "60133333333",
		"device_uuid":       devices[2].UUID,
		"verification_code": "3333",
	}

	jsonBytes3, _ := json.Marshal(userData3)

	TestHelper.Request("POST", jsonBytes3, requestURL, "")

	user1 := &User{}

	DB.Table("users").Where("guid = ?", users[0].GUID).Find(&user1)

	require.Equal(testingT{t}, 4.00, user1.Wallet)
}

func TestProfileImageSizeValidationForCreateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	sampleData.SmsHistory("register", "1111", "60111111111")

	params := map[string]string{
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       device.UUID,
		"verification_code": "1111",
	}

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/profile_image_larger.jpg", "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 413, status)
	require.Equal(testingT{t}, "File size exceeded the limit.", errorData["title"])
}

func TestProfileImageTypeValidationForCreateUser(t *testing.T) {
	DB.Exec("TRUNCATE TABLE devices;")

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	sampleData.SmsHistory("register", "1111", "60111111111")

	params := map[string]string{
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       device.UUID,
		"verification_code": "1111",
	}

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/files/test_pdf_file.pdf", "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 400, status)
	require.Equal(testingT{t}, "Invalid file type.", errorData["title"])
}

func TestCreateUserViaPhoneNumber(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	sampleData.SmsHistory("register", "1111", "60111111111")

	params := map[string]string{
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       device.UUID,
		"verification_code": "1111",
	}

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/profile_image_smaller.png", "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	accessToken := data["access_token"].(map[string]interface{})

	user := data["user"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), accessToken["expired"])
	require.NotEmpty(testingT{t}, accessToken["token"])
	require.NotEmpty(testingT{t}, user["id"])
	require.NotEmpty(testingT{t}, user["guid"])
	require.Equal(testingT{t}, params["email"], user["email"])
	require.Equal(testingT{t}, params["name"], user["name"])
	require.Equal(testingT{t}, params["phone_no"], user["phone_no"])
	require.Equal(testingT{t}, "phone_no", user["register_by"])
	require.NotEmpty(testingT{t}, user["referral_code"])
	require.Empty(testingT{t}, user["bank_account_name"])
	require.Empty(testingT{t}, user["bank_account_number"])
	require.Empty(testingT{t}, user["bank_country"])
	require.Empty(testingT{t}, user["bank_name"])

	response, _ := http.Get(user["profile_picture"].(string))

	require.Equal(testingT{t}, 200, response.StatusCode)
}

func TestCreateUserViaFacebook(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleDevice, _ := sampleData.DeviceWithoutUserGUID()

	sampleData.SmsHistory("register", "1111", "60111111111")

	params := map[string]string{
		"facebook_id":       "100013413336774",
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       sampleDevice.UUID,
		"verification_code": "1111",
	}

	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	status, _, body := TestHelper.MultipartRequest(requestURL, "POST", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/profile_image_smaller.png", "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	accessToken := data["access_token"].(map[string]interface{})

	user := data["user"].(map[string]interface{})

	device := &Device{}

	DB.Model(&Device{}).Where("guid = ?", sampleDevice.GUID).Find(&device)

	require.Equal(testingT{t}, user["guid"], *device.UserGUID)
	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), accessToken["expired"])
	require.NotEmpty(testingT{t}, accessToken["token"])
	require.NotEmpty(testingT{t}, user["id"])
	require.NotEmpty(testingT{t}, user["guid"])
	require.Equal(testingT{t}, params["facebook_id"], user["facebook_id"])
	require.Equal(testingT{t}, params["email"], user["email"])
	require.Equal(testingT{t}, params["name"], user["name"])
	require.Equal(testingT{t}, params["phone_no"], user["phone_no"])
	require.Equal(testingT{t}, "facebook", user["register_by"])
	require.NotEmpty(testingT{t}, user["referral_code"])
	require.Empty(testingT{t}, user["bank_account_name"])
	require.Empty(testingT{t}, user["bank_account_number"])
	require.Empty(testingT{t}, user["bank_country"])
	require.Empty(testingT{t}, user["bank_name"])

	response, _ := http.Get(user["profile_picture"].(string))

	require.Equal(testingT{t}, 200, response.StatusCode)
}

func TestAccessTokenRequireWhenToViewUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/8c2e6ea5-5c56-5050-ae37-a44b88e612a7", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Access token error", errorData["title"])
}

func TestViewUserDetails(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	accessToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s", TestServer.URL, users[0].GUID)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	fmt.Println(data)

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, users[0].GUID, data["guid"])
	require.Equal(testingT{t}, users[0].Name, data["name"])
	require.Equal(testingT{t}, users[0].Email, data["email"])
	require.Equal(testingT{t}, users[0].PhoneNo, data["phone_no"])
	require.Equal(testingT{t}, users[0].RegisterBy, data["register_by"])
	require.Equal(testingT{t}, users[0].ReferralCode, data["referral_code"])
}

func TestAccessTokenRequireWhenToUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/8c2e6ea5-5c56-5050-ae37-a44b88e612a7", TestServer.URL)

	status, _, body := TestHelper.Request("PATCH", []byte{}, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Access token error", errorData["title"])
}

func TestProfileImageSizeValidationForUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	accessToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s", TestServer.URL, users[0].GUID)

	params := map[string]string{
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       device.UUID,
		"verification_code": "1111",
	}

	status, _, body := TestHelper.MultipartRequest(requestURL, "PATCH", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/profile_image_larger.jpg", accessToken.Token)

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	fmt.Println(errorData)
	require.Equal(testingT{t}, 413, status)
	require.Equal(testingT{t}, "File size exceeded the limit.", errorData["title"])
}

func TestProfileImageTypeValidationForUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	accessToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s", TestServer.URL, users[0].GUID)

	params := map[string]string{
		"name":              "Test User",
		"email":             "testuser@mediacliq.my",
		"phone_no":          "60111111111",
		"device_uuid":       device.UUID,
		"verification_code": "1111",
	}

	status, _, body := TestHelper.MultipartRequest(requestURL, "PATCH", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/files/test_pdf_file.pdf", accessToken.Token)

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 400, status)
	require.Equal(testingT{t}, "Invalid file type.", errorData["title"])
}

func TestErrorAccessTokenBelongToOtherUserDuringUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	accessToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/8c2e6ea5-5c56-5050-ae37-a44b88e612a7", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken.Token)

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Your access token belong to other user", errorData["title"])
}

func TestUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	accessToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/users/%s", TestServer.URL, users[0].GUID)

	params := map[string]string{
		"name":     "Update Test User",
		"email":    "updatetestuseremail@mediacliq.my",
		"phone_no": "60199999999",
	}

	status, _, body := TestHelper.MultipartRequest(requestURL, "PATCH", params, "profile_picture",
		os.Getenv("GOPATH")+"src/bitbucket.org/cliqers/shoppermate-api/test/images/profile_image_smaller.png", accessToken.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, params["name"], data["name"])
	require.Equal(testingT{t}, params["email"], data["email"])

	response, _ := http.Get(data["profile_picture"].(string))

	require.Equal(testingT{t}, 200, response.StatusCode)
}
