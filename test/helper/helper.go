package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
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
	DB.Exec("TRUNCATE TABLE shopping_lists;")
	DB.Exec("TRUNCATE TABLE shopping_list_items;")
	DB.Exec("TRUNCATE TABLE shopping_list_item_images;")
	DB.Exec("TRUNCATE TABLE occasions;")
	DB.Exec("TRUNCATE TABLE devices;")
	DB.Exec("TRUNCATE TABLE sms_histories;")
	DB.Exec("TRUNCATE TABLE ads;")
	DB.Exec("TRUNCATE TABLE category;")
	DB.Exec("TRUNCATE TABLE subcategory;")
	DB.Exec("TRUNCATE TABLE generic;")
	DB.Exec("TRUNCATE TABLE deal_cashbacks;")
	DB.Exec("TRUNCATE TABLE deal_cashback_status;")
	DB.Exec("TRUNCATE TABLE deal_cashback_transactions;")
	DB.Exec("TRUNCATE TABLE ads_grocer")
	DB.Exec("TRUNCATE TABLE grocer;")
	DB.Exec("TRUNCATE TABLE grocer_location;")
	DB.Exec("TRUNCATE TABLE event;")
	DB.Exec("TRUNCATE TABLE event_deal;")
	DB.Exec("TRUNCATE TABLE settings;")
	DB.Exec("TRUNCATE TABLE transactions;")
	DB.Exec("TRUNCATE TABLE transaction_types;")
	DB.Exec("TRUNCATE TABLE transaction_statuses;")
	DB.Exec("TRUNCATE TABLE referral_cashback_transactions;")
	DB.Exec("TRUNCATE TABLE cashout_transactions;")

}

// UploadToAmazonS3 is a function to upload test file into amazon S3.
func (h *Helper) UploadToAmazonS3(uploadPath string, file multipart.File) (map[string]string, *systems.ErrorData) {
	// Amazon S3 Config
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION_NAME")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	fileSystem := &filesystem.FileSystem{}
	amazonS3FileSystem := fileSystem.Driver("amazonS3").(*filesystem.AmazonS3Upload)
	amazonS3FileSystem.AccessKey = accessKey
	amazonS3FileSystem.SecretKey = secretKey
	amazonS3FileSystem.Region = region
	amazonS3FileSystem.BucketName = bucketName

	localUploadPath := os.Getenv("GOPATH") + os.Getenv("STORAGE_PATH")

	uploadedFile, error := amazonS3FileSystem.Upload(file, localUploadPath, uploadPath)

	return uploadedFile, error
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

	var data interface{}

	json.Unmarshal(body.Bytes(), &data)

	return resp.StatusCode, resp.Header, data

}
