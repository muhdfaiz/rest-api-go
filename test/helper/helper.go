package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

// Helper contain useful functions can be used during testing.
type Helper struct{}

// Setup function used to initialize testing.
func (h *Helper) Setup() {
	binding.Validator = new(systems.DefaultValidator)

	h.LoadEnv()

	h.TruncateDatabase()
}

// Teardown function is a task that will be used after finish the test.
func (h *Helper) Teardown() {
	h.TruncateDatabase()
}

// LoadEnv function used to load env before start testing.
func (h *Helper) LoadEnv() {
	err := godotenv.Load(os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// TruncateDatabase function used to truncate all table inside test database.
// Useful during testing. For example, you want to clean database before start testing
// or after testing.
func (h *Helper) TruncateDatabase() {
	Database := &systems.Database{}

	DB := Database.Connect("test")

	DB.Exec("TRUNCATE TABLE users;")
	DB.Exec("TRUNCATE TABLE devices;")
	DB.Exec("TRUNCATE TABLE sms_histories;")
	DB.Exec("TRUNCATE TABLE transactions;")
	DB.Exec("TRUNCATE TABLE referral_cashback_transactions;")
}

func (h *Helper) Request(method string, jsonString []byte, url string, token string) (int, http.Header, interface{}) {
	var jsonStr = []byte(jsonString)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	resp.Header.Set("Content-Type", "application/json")

	var data interface{}

	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &data)

	return resp.StatusCode, resp.Header, data

}

func (h *Helper) MultipartRequest(uri string, method string, params map[string]string, fileParam, filePath string, token string) (int, http.Header, interface{}) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	if filePath != "" {
		file, _ := os.Open(filePath)
		defer file.Close()

		part, _ := writer.CreateFormFile(fileParam, filepath.Base(filePath))
		io.Copy(part, file)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.Close()

	req, err := http.NewRequest(method, uri, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body = &bytes.Buffer{}

	_, err = body.ReadFrom(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(body)

	var data interface{}

	json.Unmarshal(body.Bytes(), &data)

	return resp.StatusCode, resp.Header, data

}
