package v1_1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

func TestErrorRequiredFieldDuringCreateDevice(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/devices", TestServer.URL)

	params := map[string]string{
		"user_guid":   "8c2e6ea5-5c56-5050-ae37-a44b88e612a7",
		"os":          "",
		"model":       "",
		"uuid":        "",
		"push_token":  "",
		"app_version": "",
	}

	jsonBytes, _ := json.Marshal(params)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	error := body.(map[string]interface{})["errors"].(map[string]interface{})
	errorDetail := error["detail"].(map[string]interface{})

	require.Equal(testingT{t}, 422, status)
	require.Equal(testingT{t}, "Validation failed.", error["title"])
	require.NotEmpty(testingT{t}, errorDetail["os"])
	require.NotEmpty(testingT{t}, errorDetail["model"])
	require.NotEmpty(testingT{t}, errorDetail["uuid"])
	require.NotEmpty(testingT{t}, errorDetail["push_token"])
	require.NotEmpty(testingT{t}, errorDetail["app_version"])
}

func TestErrorUserGUIDNotExistDuringCreateDevice(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/devices", TestServer.URL)

	params := map[string]string{
		"user_guid":   "8c2e6ea5-5c56-5050-ae37-a44b88e612a7",
		"os":          "Android",
		"model":       "Xiaomi MI 3W",
		"uuid":        "F4EE037F5340586BE625EE7678797DD0",
		"push_token":  "cuh2h6UWq-0:APA91bHKfVGhWAoO_x2xAZ2WfdA-Xraw78Eyp58fMsqeJ_jbrt1vkZkTnwKPFvgYLtVP55nXBKseqpo0mqQ7JmK3iPGYFjptQj72HizU4FJdqbQ2AEgvnNxk6gAPeT9eRv9FNQb_TJyB",
		"app_version": "0.9.9.3-STAGING-BETA",
	}

	jsonBytes, _ := json.Marshal(params)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	error := body.(map[string]interface{})["errors"].(map[string]interface{})
	errorDetail := error["detail"].(map[string]interface{})

	require.Equal(testingT{t}, 404, status)
	require.Equal(testingT{t}, "User not exists.", error["title"])
	require.NotEmpty(testingT{t}, errorDetail["guid"])
}

func TestErrorDuplicateDeviceDuringCreateDevice(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	device, _ := sampleData.DeviceWithoutUserGUID()

	requestURL := fmt.Sprintf("%s/v1_1/devices", TestServer.URL)

	params := map[string]string{
		"os":          "Android",
		"model":       "Xiaomi MI 3W",
		"uuid":        device.UUID,
		"push_token":  "cuh2h6UWq-0:APA91bHKfVGhWAoO_x2xAZ2WfdA-Xraw78Eyp58fMsqeJ_jbrt1vkZkTnwKPFvgYLtVP55nXBKseqpo0mqQ7JmK3iPGYFjptQj72HizU4FJdqbQ2AEgvnNxk6gAPeT9eRv9FNQb_TJyB",
		"app_version": "0.9.9.3-STAGING-BETA",
	}

	jsonBytes, _ := json.Marshal(params)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	error := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 409, status)
	require.Equal(testingT{t}, "Device already exists.", error["title"])
}

func TestSuccessCreateDeviceWithUserGUID(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	params := map[string]string{
		"user_guid":   users[0].GUID,
		"os":          "Android",
		"model":       "Xiaomi MI 3W",
		"uuid":        "F4EE037F5340586BE625EE7678797DD0",
		"push_token":  "cuh2h6UWq-0:APA91bHKfVGhWAoO_x2xAZ2WfdA-Xraw78Eyp58fMsqeJ_jbrt1vkZkTnwKPFvgYLtVP55nXBKseqpo0mqQ7JmK3iPGYFjptQj72HizU4FJdqbQ2AEgvnNxk6gAPeT9eRv9FNQb_TJyB",
		"app_version": "0.9.9.3-STAGING-BETA",
	}

	jsonBytes, _ := json.Marshal(params)

	requestURL := fmt.Sprintf("%s/v1_1/devices", TestServer.URL)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	fmt.Println(body)
	response := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, params["user_guid"], response["user_guid"])
	require.Equal(testingT{t}, params["os"], response["os"])
	require.Equal(testingT{t}, params["model"], response["model"])
	require.Equal(testingT{t}, params["uuid"], response["uuid"])
	require.Equal(testingT{t}, params["push_token"], response["push_token"])
	require.Equal(testingT{t}, params["app_version"], response["app_version"])

}

func TestSuccessCreateDeviceWithoutUserGUID(t *testing.T) {
	TestHelper.TruncateDatabase()

	params := map[string]string{
		"os":          "Android",
		"model":       "Xiaomi MI 3W",
		"uuid":        "F4EE037F5340586BE625EE7678797DD0",
		"push_token":  "cuh2h6UWq-0:APA91bHKfVGhWAoO_x2xAZ2WfdA-Xraw78Eyp58fMsqeJ_jbrt1vkZkTnwKPFvgYLtVP55nXBKseqpo0mqQ7JmK3iPGYFjptQj72HizU4FJdqbQ2AEgvnNxk6gAPeT9eRv9FNQb_TJyB",
		"app_version": "0.9.9.3-STAGING-BETA",
	}

	jsonBytes, _ := json.Marshal(params)

	requestURL := fmt.Sprintf("%s/v1_1/devices", TestServer.URL)

	status, _, body := TestHelper.Request("POST", jsonBytes, requestURL, "")

	fmt.Println(body)
	response := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.Empty(testingT{t}, response["user_guid"])
	require.Equal(testingT{t}, params["os"], response["os"])
	require.Equal(testingT{t}, params["model"], response["model"])
	require.Equal(testingT{t}, params["uuid"], response["uuid"])
	require.Equal(testingT{t}, params["push_token"], response["push_token"])
	require.Equal(testingT{t}, params["app_version"], response["app_version"])

}

func TestErrorAccessTokenDuringDeleteDevice(t *testing.T) {
	requestURL := fmt.Sprintf("%s/v1_1/devices/%s", TestServer.URL, "F4EE037F5340586BE625EE7678797DD0")

	status, _, body := TestHelper.Request("DELETE", []byte{}, requestURL, "")

	error := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 401, status)
	require.Equal(testingT{t}, "Access token error", error["title"])
}

func TestErrorDeviceUUIDNotExistDuringDeleteDevice(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/devices/%s", TestServer.URL, "F4EE037F5340586BE625EE7678792222")

	status, _, body := TestHelper.Request("DELETE", []byte{}, requestURL, jwtToken.Token)

	error := body.(map[string]interface{})["errors"].(map[string]interface{})

	require.Equal(testingT{t}, 404, status)
	require.Equal(testingT{t}, "Device not exists.", error["title"])
}

func TestSuccessfulSoftDeleteDevice(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/devices/%s", TestServer.URL, device.UUID)

	status, _, _ := TestHelper.Request("DELETE", []byte{}, requestURL, jwtToken.Token)

	require.Equal(testingT{t}, 200, status)

	deletedDevice := &Device{}

	DB.Unscoped().Model(&Device{}).Where(&Device{UUID: device.UUID}).Find(&deletedDevice)

	require.NotEmpty(testingT{t}, deletedDevice.DeletedAt)
}

func TestErrorDeviceNotExistDuringUpdateDevice(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/devices/%s", TestServer.URL, "F4EE037F5340586BE625EE7678797111")

	params := map[string]string{
		"user_guid":   users[1].GUID,
		"os":          "iOS",
		"model":       "Iphone6s",
		"uuid":        "F4EE037F5340586BE625EE7678797DD0",
		"app_version": "1.0-STAGING",
	}

	jsonBytes, _ := json.Marshal(params)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwtToken.Token)

	error := body.(map[string]interface{})["errors"].(map[string]interface{})
	require.Equal(testingT{t}, 404, status)
	require.Equal(testingT{t}, "Device not exists.", error["title"])
}

func TestSuccessfulUpdateDevice(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	users := sampleData.Users()

	device := sampleData.DeviceWithUserGUID(users[0].GUID)

	jwt := &systems.Jwt{}

	jwtToken, _ := jwt.GenerateToken(users[0].GUID, users[0].PhoneNo, device.UUID, "")

	requestURL := fmt.Sprintf("%s/v1_1/devices/%s", TestServer.URL, device.UUID)

	params := map[string]string{
		"user_guid":   users[1].GUID,
		"os":          "iOS",
		"model":       "Iphone6s",
		"app_version": "1.0-STAGING",
	}

	jsonBytes, _ := json.Marshal(params)

	status, _, body := TestHelper.Request("PATCH", jsonBytes, requestURL, jwtToken.Token)

	response := body.(map[string]interface{})["data"].(map[string]interface{})

	require.Equal(testingT{t}, 200, status)
	require.Equal(testingT{t}, params["user_guid"], response["user_guid"])
	require.Equal(testingT{t}, params["os"], response["os"])
	require.Equal(testingT{t}, params["model"], response["model"])
	require.Equal(testingT{t}, params["app_version"], response["app_version"])
}
