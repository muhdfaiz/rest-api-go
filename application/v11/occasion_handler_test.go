package v11

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestViewAllActiveOccasionsShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	sampleData.Occasions()

	requestURL := fmt.Sprintf("%s/v1_1/shopping_lists/occasions", TestServer.URL)

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	data := body.(map[string]interface{})["data"].([]interface{})

	occasion1 := data[0].(map[string]interface{})
	occasion2 := data[1].(map[string]interface{})

	assert.Equal(t, 200, status)
	assert.Len(t, data, 2)
	assert.Equal(t, "Field Trip", occasion1["name"])
	assert.Equal(t, "field_trip", occasion1["slug"])
	assert.Equal(t, 1.00, occasion1["active"])

	assert.Equal(t, "Travel", occasion2["name"])
	assert.Equal(t, "travel", occasion2["slug"])
	assert.Equal(t, 1.00, occasion2["active"])
}

func TestViewLatestActiveOccasionsShouldSuccess(t *testing.T) {
	TestHelper.TruncateDatabase()

	sampleData := SampleData{DB: DB}

	occasions := sampleData.Occasions()

	DB.Model(&Occasion{}).Where(&Occasion{GUID: occasions[0].GUID}).UpdateColumn("updated_at", time.Now().UTC().Add(time.Hour*24*7))

	requestURL := fmt.Sprintf("%s/v1_1/shopping_lists/occasions?last_sync_date=%s", TestServer.URL, time.Now().UTC().Add(time.Hour*24*3).Format(time.RFC3339))

	status, _, body := TestHelper.Request("GET", []byte{}, requestURL, "")

	data := body.(map[string]interface{})["data"].([]interface{})

	occasion1 := data[0].(map[string]interface{})

	assert.Equal(t, 200, status)
	assert.Len(t, data, 1)
	assert.Equal(t, "Field Trip", occasion1["name"])
	assert.Equal(t, "field_trip", occasion1["slug"])
	assert.Equal(t, 1.00, occasion1["active"])
}
