package v1_1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewAllSettingShouldBeSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Settings("true", "2", "2")

	requestURL := fmt.Sprintf("%s/v1_1/settings", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	settings := body.(map[string]interface{})["data"].([]interface{})

	assert.Equal(t, 200, status)
	assert.Equal(t, 3, len(settings))
	assert.Equal(t, "referral_active", settings[0].(map[string]interface{})["slug"])
	assert.Equal(t, "referral_price", settings[1].(map[string]interface{})["slug"])
	assert.Equal(t, "max_referral_user", settings[2].(map[string]interface{})["slug"])
}
