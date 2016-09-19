package filesystem

// import (
// 	"bytes"
// 	"io"
// 	"io/ioutil"
// 	"os"
// 	"path"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"bitbucket.org/shoppermate/systems"

// 	filetype "gopkg.in/h2non/filetype.v0"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	uuid "github.com/satori/go.uuid"
// )

// // Initialize object needed
// var (
// 	Config    = &systems.Configs{}
// 	ErrorMesg = &systems.Error{}
// 	Helper    = &systems.Helpers{}
// )

// type StoragePath interface {
// 	Get() string
// }

// type Storage struct{}

// type AmazonS3Interface interface {
// 	SetBucketName() string
// 	SetUploadPath() string
// }

// type AmazonS3Upload struct{}

// // UploadService will hold all variable and function for upload file
// type UploadService struct {
// 	File                 io.Reader
// 	FileAttribute        string
// 	FileName             string
// 	MaximumFileSizeAllow int64 //size must be in bytes
// 	FileTypeAllow        []string
// 	FileSize             int64
// 	FileType             string
// 	FilePath             string
// }

// // Get function used to retrieve storage path for file
// func (s *Storage) Get() string {
// 	return os.Getenv("GOPATH") + "src/bitbucket.org/shoppermate/storage/"
// }

// // Get function used to retrieve bucket name for Amazon S3
// func (a *AmazonS3) Get() string {
// 	return Config.Get("app.yaml", "aws_bucket_name", "shoppermate")
// }

// // UploadToAmazonS3Storage is a function to upload file to Amazon S3.
// func (us *UploadService) UploadToAmazonS3Storage(sp StoragePath, bn BucketName) (interface{}, *systems.ErrorData) {

// 	// Upload to local storage first before upload to Amazon S3
// 	_, err := us.UploadToLocalStorage(sp, true)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Retrieve AWS Config
// 	creds := credentials.NewStaticCredentials(Config.Get("app.yaml", "aws_access_key_id", ""), Config.Get("app.yaml", "aws_secret_access_key", ""), "")

// 	_, err1 := creds.Get()
// 	if err1 != nil {
// 		return nil, ErrorMesg.InternalServerError(err1.Error(), systems.ErrorAmazonService)
// 	}

// 	// Initialize AWS Session with credentials
// 	awsSession := s3.New(session.New(), aws.NewConfig().WithRegion(Config.Get("app.yaml", "aws_region_name", "")).WithCredentials(creds))

// 	// Read file uploaded in local storage
// 	file, err2 := os.Open(sp.Get() + us.FileName)
// 	if err2 != nil {
// 		return nil, ErrorMesg.InternalServerError(err, systems.CannotReadFile)
// 	}
// 	defer file.Close()

// 	// Read file content to buffer
// 	buffer := make([]byte, us.FileSize)

// 	file.Read(buffer)
// 	fileBytes := bytes.NewReader(buffer)

// 	// Specify path or folder that will be used to upload image
// 	path := "/profile_image/" + us.FileName
// 	params := &s3.PutObjectInput{
// 		Bucket:        aws.String(bn.Get()),
// 		Key:           aws.String(path),
// 		Body:          fileBytes,
// 		ContentLength: aws.Int64(us.FileSize),
// 		ContentType:   aws.String(us.FileType),
// 	}

// 	_, err3 := awsSession.PutObject(params)
// 	if err3 != nil {
// 		return nil, ErrorMesg.InternalServerError(err3.Error(), systems.ErrorAmazonService)
// 	}

// 	return us.setUploadResult(), nil
// }

// // UploadToLocalStorage is a function to upload file in the local storage.
// // Local storage path is define app config.
// // Return filename, filesize & filetype
// func (us *UploadService) UploadToLocalStorage(sp StoragePath, deleteAfterUploaded bool) (map[string]interface{}, *systems.ErrorData) {
// 	if us.File == nil || us.FileAttribute == "" {
// 		return nil, ErrorMesg.GenericError("400", systems.CannotReadFile, systems.TitleFileEmptyError, "", systems.ErrorFileEmpty)
// 	}

// 	// Generate Unique File Name
// 	newFileName := us.GenerateUniqueFileName() + path.Ext(us.FileName)

// 	// Create new file in local storage
// 	localFile, err := os.Create(sp.Get() + newFileName)

// 	if err != nil {
// 		return nil, ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
// 	}
// 	defer localFile.Close()

// 	// Copy file from request to created file in local storage
// 	_, err = io.Copy(localFile, us.File)
// 	if err != nil {
// 		return nil, ErrorMesg.InternalServerError(err.Error(), systems.CannotCopyFile)
// 	}

// 	// Retrieve file information
// 	fileInfo, err := localFile.Stat()
// 	if err != nil {
// 		return nil, ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
// 	}

// 	// Set filesize to global variable in UploadService object
// 	us.FileSize = fileInfo.Size()

// 	// Validate File size
// 	err1 := us.ValidateFileSize()
// 	if err1 != nil {
// 		return nil, err1
// 	}

// 	// Detect & Validate File type
// 	err1 = us.ValidateFileType()
// 	if err1 != nil {
// 		return nil, err1
// 	}

// 	err = os.Remove(us.FilePath)

// 	return us.setUploadResult(), nil
// }

// // ValidateFileType function used to validate file type allow to upload
// func (us *UploadService) ValidateFileType() *systems.ErrorData {
// 	// Skip validation file type if FileTypeAllow empty
// 	if us.FileTypeAllow == nil {
// 		return nil
// 	}

// 	// Detect filetype
// 	err := us.DetectFileType()
// 	if err != nil {
// 		return err
// 	}

// 	for _, allowFileType := range us.FileTypeAllow {
// 		if allowFileType == us.FileType {
// 			return nil
// 		}
// 	}

// 	return ErrorMesg.InvalidFileTypeError(strings.Join(us.FileTypeAllow, ","))
// }

// // DetectFileType function will validate file type allow to upload or not
// func (us *UploadService) DetectFileType() *systems.ErrorData {
// 	if us.FilePath == "" || us.FileName == "" {
// 		return ErrorMesg.InternalServerError(systems.ErrorFileEmpty, systems.CannotReadFile)
// 	}

// 	buffer, err := ioutil.ReadFile(us.FilePath + us.FileName)
// 	if err != nil {
// 		return ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
// 	}

// 	filetype, err := filetype.Match(buffer)
// 	if err != nil {
// 		return ErrorMesg.InternalServerError(err.Error(), systems.CannotDetectFileType)
// 	}
// 	us.FileType = filetype.Extension
// 	return nil

// }

// // ValidateFileSize function used to verify if file size want to upload is not bigger than system allowed
// func (us *UploadService) ValidateFileSize() *systems.ErrorData {
// 	if us.MaximumFileSizeAllow == 0 {
// 		return nil
// 	}

// 	if us.FileSize >= us.MaximumFileSizeAllow {
// 		return ErrorMesg.FileSizeExceededLimit(us.FileAttribute, strconv.FormatInt(us.MaximumFileSizeAllow/1000, 10))
// 	}
// 	return nil
// }

// // GenerateUniqueFileName function used to generate unique file name using uuid v5 with timestamp
// func (us *UploadService) GenerateUniqueFileName() string {
// 	namespace := uuid.NamespaceDNS
// 	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
// 	return uuid.NewV3(namespace, timestamp).String()
// }

// // setUploadResult function used format output after successfully upload file
// func (us *UploadService) setUploadResult() map[string]interface{} {
// 	uploadResult := make(map[string]interface{})
// 	uploadResult["filename"] = us.FileName
// 	uploadResult["filetype"] = us.FileType
// 	uploadResult["filesize"] = us.FileSize
// 	uploadResult["filepath"] = us.FileName

// 	return uploadResult
// }
