package v1_1

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoginViaPhoneShouldReturnValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/phone", TestServer.URL)

	postData := LoginViaPhone{
		PhoneNo: "",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["phone_no"])
}

func TestLoginViaPhoneShouldReturnPhoneNumberNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/phone", TestServer.URL)

	postData := LoginViaPhone{
		PhoneNo: "601111111111",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "User not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["phone_no"])
}

func TestLoginViaPhoneUsingDebugModeShouldNotSentSms(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	sampleData.DeviceWithUserGUID(users[0].GUID)

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/phone?debug=1", TestServer.URL)

	postData := LoginViaPhone{
		PhoneNo: users[0].PhoneNo,
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, users[0].GUID, data["user_guid"])

	smsHistory := &SmsHistory{}

	DB.Model(&SmsHistory{}).First(&smsHistory)

	require.Empty(t, smsHistory.GUID)
}

func TestLoginViaPhoneUsingDebugModeShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("60174862127", "muhdfaiz@mediacliq.my")

	sampleData.DeviceWithUserGUID(user.GUID)

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/phone", TestServer.URL)

	postData := LoginViaPhone{
		PhoneNo: user.PhoneNo,
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, user.GUID, data["user_guid"])

	smsHistory := &SmsHistory{}

	DB.Model(&SmsHistory{}).First(&smsHistory)

	require.Equal(t, user.PhoneNo, smsHistory.RecipientNo)
	require.Equal(t, "login", smsHistory.Event)
	require.NotEmpty(t, smsHistory.VerificationCode)
}

func TestLoginViaFacebookShouldReturnValidationError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/facebook", TestServer.URL)

	postData := LoginViaFacebook{
		FacebookID: "",
		DeviceUUID: "",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 422, status)
	require.Equal(t, 422, status)
	require.Equal(t, "Validation failed.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["device_uuid"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["facebook_id"])

}

func TestLoginViaFacebookShouldReturnFacebookIDNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/facebook", TestServer.URL)

	postData := LoginViaFacebook{
		FacebookID: "1231232343434",
		DeviceUUID: "1231231231231",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "User not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["facebook_id"])
}

func TestLoginViaFacebookShouldReturnDeviceNotExistError(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	user := sampleData.User("60174862127", "muhdfaiz@mediacliq.my")

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/facebook", TestServer.URL)

	postData := LoginViaFacebook{
		FacebookID: *user.FacebookID,
		DeviceUUID: "1231231231231",
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 404, status)
	require.Equal(t, "Device not exists.", errors["title"])
	require.NotEmpty(t, errors["detail"].(map[string]interface{})["uuid"])
}

func TestLoginViaFacebookShouldUpdateDeviceUserGUIDWhenDeviceUserGUIDEmpty(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device, _ := sampleData.DeviceWithoutUserGUID()

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/facebook", TestServer.URL)

	facebookID := users[0].FacebookID

	postData := LoginViaFacebook{
		DeviceUUID: device.UUID,
		FacebookID: *facebookID,
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, _ := TestHelper.Request("POST", jsonBytes, requestURL, "")

	updatedDevice := &Device{}

	DB.Model(&Device{}).Where(&Device{GUID: device.GUID}).Find(updatedDevice)

	require.Equal(t, 200, status)
	require.Equal(t, users[0].GUID, *updatedDevice.UserGUID)
}

func TestLoginViaFacebookShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	DB.Model(&Device{}).Where(&Device{GUID: device.GUID}).Select("deleted_at").Update(map[string]interface{}{"deleted_at": time.Now().UTC()})

	requestURL := fmt.Sprintf("%s/v1_1/auth/login/facebook", TestServer.URL)

	facebookID := users[0].FacebookID

	postData := LoginViaFacebook{
		DeviceUUID: device.UUID,
		FacebookID: *facebookID,
	}

	jsonBytes, _ := json.Marshal(postData)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	updatedDevice := &Device{}

	DB.Model(&Device{}).Where(&Device{GUID: device.GUID}).Find(updatedDevice)

	require.Equal(t, 200, status)
	require.Nil(t, updatedDevice.DeletedAt)

	require.Equal(t, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), data["access_token"].(map[string]interface{})["expired"])
	require.Equal(t, users[0].GUID, data["user"].(map[string]interface{})["guid"])
}

func TestLogoutShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	requestURL := fmt.Sprintf("%s/v1_1/auth/logout", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestLogoutShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/auth/logout", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.Equal(t, "Successfully logout", data["message"])

	deletedDevice := &Device{}

	userGUID := users[0].GUID

	DB.Unscoped().Model(&Device{}).Where(&Device{UserGUID: &userGUID}).Find(deletedDevice)

	require.NotNil(t, deletedDevice.DeletedAt)
}

func TestRefreshAccessTokenShouldReturnAccessTokenError(t *testing.T) {
	TestHelper.TruncateDatabase()

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZV9ubyI6IjYwMTc0ODYyMTI3IiwiYXVkIjoiOGMyZTZlYTUtNWM1Ni01MDUwLWFlMzctYTQ0Yjg4ZTYxMmE3IiwiZXhwIjoxNDg3NjcwNjA3LCJqdGkiOiJGNEVFMDM3RjUzNDA1ODZCRTYyNUVFNzY3ODc5N0REMCIsImlhdCI6MTQ4NzA2NTgwNywiaXNzIjoiaHR0cDovL2FwaS5zaG9wcGVybWF0ZS1hcGkuY29tIiwibmJmIjoxNDg3MDY1ODA3LCJzdWIiOiI4YzJlNmVhNS01YzU2LTUwNTAtYWUzNy1hNDRiODhlNjEyYTcifQ.71ZzAnZELFTnsnh8wRCDyG4IKzOaSv3VJDxYnHk6GHY"

	requestURL := fmt.Sprintf("%s/v1_1/auth/refresh", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, accessToken)

	errors := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(t, 401, status)
	require.Equal(t, "Access token error", errors["title"])
}

func TestRefreshAccessTokenShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt, _ := JWT.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/auth/refresh", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, jwt.Token)

	data := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(t, 200, status)
	require.NotEmpty(t, data["token"])
	require.Equal(t, time.Now().UTC().AddDate(0, 0, 7).Format(time.RFC3339), data["expired"])
}
