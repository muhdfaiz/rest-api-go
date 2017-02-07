package v1_1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"bitbucket.org/cliqers/shoppermate-api/test/helper"
)

var (
	Router     = gin.Default()
	TestHelper = helper.Helper{}
	TestServer *httptest.Server
	DB         *gorm.DB
)

func TestMain(m *testing.M) {
	TestHelper.Setup()

	DB = Database.Connect("test")

	TestServer = httptest.NewServer(InitializeObjectAndSetRoutesV1_1(Router, DB))

	ret := m.Run()

	if ret == 0 {
		//TestHelper.Teardown()
		TestServer.Close()
	}

	os.Exit(ret)
}

// TestNumericFieldsDuringCreateUser function used to test if API will return error or not if the input
// for numeric field like facebook_id not numeric.
func TestNumericFieldsDuringCreateUser(t *testing.T) {
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

	assert.Equal(t, 422, status, "Status code must be 422")
	assert.NotEmpty(t, errorDetail["facebook_id"])
	assert.NotEmpty(t, errorDetail["phone_no"])
	assert.NotEmpty(t, errorDetail["verification_code"])
}

// TestRequiredFieldsDuringCreateUser function used to test if API will return error if the input data not
// contain all the required fields.
func TestRequiredFieldsDuringCreateUser(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/users", TestServer.URL)

	userData := map[string]string{
		"facebook_id": "123123123123",
	}

	jsonBytes, _ := json.Marshal(userData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})
	errorDetail := errorData["detail"].(map[string]interface{})

	assert.Equal(t, 422, status, "Status code must be 422")
	assert.NotEmpty(t, errorDetail["name"])
	assert.NotEmpty(t, errorDetail["email"])
	assert.NotEmpty(t, errorDetail["phone_no"])
	assert.NotEmpty(t, errorDetail["verification_code"])
	assert.NotEmpty(t, errorDetail["device_uuid"])
}

// TestPhoneNumberMustUnique function used to test if API return error or not when the input data
// contain the same phone number that already exist in database.
func TestPhoneNumberMustUniqueDuringCreateUser(t *testing.T) {
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

	assert.Equal(t, 409, status, "Status code must be 409")
	assert.Equal(t, "Duplicate entry 'phone_no' for key '"+userData["phone_no"]+"'.", errorData["detail"].(map[string]interface{})["message"])
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

	assert.Equal(t, 400, status, "Status code must be 409")
	assert.NotEmpty(t, errorData["detail"].(map[string]interface{})["facebook_id"])
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

	assert.Equal(t, 409, status, "Status code must be 409")
	assert.Equal(t, "Duplicate entry 'facebook_id' for key '100013413336774'.", errorData["detail"].(map[string]interface{})["message"])
}

func TestReferralCodeMustValidDuringCreateUser(t *testing.T) {
	DB.Model(&Setting{}).Where("slug = ?", "referral_active").Update("value", "true")

	sampleData := SampleData{DB: DB}

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

	assert.Equal(t, 404, status, "Status code must be 404")
	assert.NotEmpty(t, errorData["detail"].(map[string]interface{})["referral_code"])
}

func TestCreateUserWithReferralCodeNotActiveAndDeviceUUIDMustValid(t *testing.T) {
	DB.Model(&Setting{}).Where("slug = ?", "referral_active").Update("value", "false")

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

	assert.Equal(t, 404, status, "Status code must be 404")
	assert.NotEmpty(t, errorData["detail"].(map[string]interface{})["uuid"])
}

// TestDebug function used to test if API can create new user and generate access token without valid verification code.
func TestCreateUserWithDebugMode(t *testing.T) {
	DB.Exec("TRUNCATE TABLE devices;")

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

	fmt.Println(data)
	fmt.Println(status)

	assert.NotEmpty(t, data["access_token"])
	assert.NotEmpty(t, data["user"])
}

// TestDebugToken function used to test if API will generate access token based on debug token value.
func TestCreateUserWithDebugToken(t *testing.T) {
	DB.Exec("TRUNCATE TABLE users;")
	DB.Exec("TRUNCATE TABLE devices;")

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

	assert.Equal(t, 200, status)
	assert.Equal(t, time.Now().UTC().Add(minutes).Format(time.RFC3339), data["access_token"].(map[string]interface{})["expired"])
	assert.NotEmpty(t, data["user"])
}

func TestCreateUserWithReferral(t *testing.T) {
	DB.Model(&Setting{}).Where("slug = ?", "referral_active").Update("value", "true")

	referralPriceSetting := &Setting{}
	DB.Model(&Setting{}).Where("slug = ?", "referral_price").Find(&referralPriceSetting)

	DB.Exec("TRUNCATE TABLE users;")
	DB.Exec("TRUNCATE TABLE devices;")

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	users := sampleData.Users()

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

	assert.Equal(t, referralPriceSetting.Value, strconv.FormatFloat(user1.Wallet, 'f', 0, 64))

	referralCashbackTransaction := &ReferralCashbackTransaction{}

	DB.Model(&ReferralCashbackTransaction{}).Where("user_guid = ?", user1.GUID).Find(&referralCashbackTransaction)

	assert.Equal(t, 200, status)
	assert.Equal(t, user1.GUID, referralCashbackTransaction.UserGUID)
	assert.Equal(t, newUser["guid"], referralCashbackTransaction.ReferrerGUID)
	assert.NotEmpty(t, referralCashbackTransaction.TransactionGUID)

	transaction := &Transaction{}

	DB.Model(&Transaction{}).Where("guid = ?", referralCashbackTransaction.TransactionGUID).Find(&transaction)

	assert.NotEmpty(t, transaction.GUID)
	assert.Equal(t, user1.GUID, transaction.UserGUID)
	assert.Equal(t, "a606113b-fb22-59f3-876f-dd05da7befc7", transaction.TransactionTypeGUID)
	assert.Equal(t, "669e13c0-eaea-5aef-a25f-6ba54b529e33", transaction.TransactionStatusGUID)
	assert.NotEmpty(t, transaction.ReferenceID)
	assert.Equal(t, referralPriceSetting.Value, strconv.FormatFloat(transaction.TotalAmount, 'f', 0, 64))

}

func TestMaxReferralPerUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	DB.Model(&Setting{}).Where("slug = ?", "referral_active").Update("value", "true")

	DB.Model(&Setting{}).Select("value").Where("slug = ?", "max_referral_user").Updates(map[string]interface{}{"value": "2"})

	referralPriceSetting := &Setting{}

	DB.Model(&Setting{}).Where("slug = ?", "referral_price").Find(&referralPriceSetting)

	sampleData := SampleData{DB: DB}

	devices, _ := sampleData.DevicesWithoutUserGUID()

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

	assert.Equal(t, "4", strconv.FormatFloat(user1.Wallet, 'f', 0, 64))
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

	assert.Equal(t, 413, status)
	assert.Equal(t, "File size exceeded the limit.", errorData["title"])
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

	assert.Equal(t, 400, status)
	assert.Equal(t, "Invalid file type.", errorData["title"])
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

	assert.Equal(t, 200, status)
	assert.Equal(t, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), accessToken["expired"])
	assert.NotEmpty(t, accessToken["token"])
	assert.NotEmpty(t, user["id"])
	assert.NotEmpty(t, user["guid"])
	assert.Equal(t, params["email"], user["email"])
	assert.Equal(t, params["name"], user["name"])
	assert.Equal(t, params["phone_no"], user["phone_no"])
	assert.Equal(t, "phone_no", user["register_by"])
	assert.NotEmpty(t, user["referral_code"])
	assert.Empty(t, user["bank_account_name"])
	assert.Empty(t, user["bank_account_number"])
	assert.Empty(t, user["bank_country"])
	assert.Empty(t, user["bank_name"])

	response, _ := http.Get(user["profile_picture"].(string))

	assert.Equal(t, 200, response.StatusCode)
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

	assert.Equal(t, user["guid"], *device.UserGUID)
	assert.Equal(t, 200, status)
	assert.Equal(t, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), accessToken["expired"])
	assert.NotEmpty(t, accessToken["token"])
	assert.NotEmpty(t, user["id"])
	assert.NotEmpty(t, user["guid"])
	assert.Equal(t, params["facebook_id"], user["facebook_id"])
	assert.Equal(t, params["email"], user["email"])
	assert.Equal(t, params["name"], user["name"])
	assert.Equal(t, params["phone_no"], user["phone_no"])
	assert.Equal(t, "facebook", user["register_by"])
	assert.NotEmpty(t, user["referral_code"])
	assert.Empty(t, user["bank_account_name"])
	assert.Empty(t, user["bank_account_number"])
	assert.Empty(t, user["bank_country"])
	assert.Empty(t, user["bank_name"])

	response, _ := http.Get(user["profile_picture"].(string))

	assert.Equal(t, 200, response.StatusCode)
}

func TestAccessTokenRequireWhenToViewUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/8c2e6ea5-5c56-5050-ae37-a44b88e612a7", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errorData["title"])
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

	assert.Equal(t, 200, status)
	assert.Equal(t, users[0].GUID, data["guid"])
	assert.Equal(t, users[0].Name, data["name"])
	assert.Equal(t, users[0].Email, data["email"])
	assert.Equal(t, users[0].PhoneNo, data["phone_no"])
	assert.Equal(t, users[0].RegisterBy, data["register_by"])
	assert.Equal(t, users[0].ReferralCode, data["referral_code"])
}

func TestAccessTokenRequireWhenToUpdateUser(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/users/8c2e6ea5-5c56-5050-ae37-a44b88e612a7", TestServer.URL)

	status, _, body := TestHelper.Request("PATCH", []byte{}, requestURL, "")

	errorData := body.(map[string]interface{})["errors"].(map[string]interface{})

	assert.Equal(t, 401, status)
	assert.Equal(t, "Access token error", errorData["title"])
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
	assert.Equal(t, 413, status)
	assert.Equal(t, "File size exceeded the limit.", errorData["title"])
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

	assert.Equal(t, 400, status)
	assert.Equal(t, "Invalid file type.", errorData["title"])
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

	assert.Equal(t, 401, status)
	assert.Equal(t, "Your access token belong to other user", errorData["title"])
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

	assert.Equal(t, 200, status)
	assert.Equal(t, params["name"], data["name"])
	assert.Equal(t, params["email"], data["email"])

	response, _ := http.Get(data["profile_picture"].(string))

	assert.Equal(t, 200, response.StatusCode)
}
