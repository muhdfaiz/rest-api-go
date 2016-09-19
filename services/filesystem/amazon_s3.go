package filesystem

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"

	filetype "gopkg.in/h2non/filetype.v0"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"bitbucket.org/shoppermate/systems"
)

var (
	Config = &systems.Configs{}
)

type AmazonS3ServiceUpload struct {
	AccessKey string
	SecretKey string
	Region    string
}

type S3UploadConfigInterface interface {
	SetAmazonS3UploadPath() string
	SetLocalUploadPath() string
	SetBucketName() string
}

// Upload function used to store file in local storage
func (asu *AmazonS3ServiceUpload) Upload(s3uploadConfig S3UploadConfigInterface, file multipart.File) (map[string]string, *systems.ErrorData) {

	localUploadService := LocalUpload{}
	uploadedFile, err := localUploadService.Upload(s3uploadConfig, file)
	if err != nil {
		return nil, err
	}

	// Read file
	localFile, err1 := os.Open(s3uploadConfig.SetLocalUploadPath() + uploadedFile["name"])
	if err1 != nil {
		return nil, ErrorMesg.InternalServerError(err1.Error(), systems.CannotReadFile)
	}
	defer localFile.Close()

	// Retrieve file size
	fileInfo, _ := localFile.Stat()
	fileSize := fileInfo.Size()

	// Read file content to buffer
	buffer := make([]byte, fileSize)

	localFile.Read(buffer)
	fileBytes := bytes.NewReader(buffer)

	// Specify path or folder that will be used to upload image
	path := s3uploadConfig.SetAmazonS3UploadPath() + uploadedFile["name"]
	params := &s3.PutObjectInput{
		Bucket:        aws.String(s3uploadConfig.SetBucketName()),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(uploadedFile["type"]),
	}

	awsSession, err := asu.CreateSession()
	if err != nil {
		return nil, err
	}

	_, err2 := awsSession.PutObject(params)
	if err2 != nil {
		return nil, ErrorMesg.InternalServerError(err2.Error(), systems.ErrorAmazonService)
	}

	// Remove Local File
	err2 = os.Remove(uploadedFile["path"])
	if err2 != nil {
		return nil, ErrorMesg.InternalServerError(err2.Error(), systems.CannotDeleteFile)
	}

	uploadedFile["path"] = fmt.Sprintf("https://s3-%s.amazonaws.com/%s%s%s", Config.Get("app.yaml", "aws_region_name", "ap-southeast-1"),
		Config.Get("app.yaml", "aws_bucket_name", ""), s3uploadConfig.SetAmazonS3UploadPath(), uploadedFile["name"])

	return uploadedFile, nil
}

func (asu *AmazonS3ServiceUpload) Delete(filepath string) {

}

// GetFileType function will validate file type allow to upload or not
func (asu *AmazonS3ServiceUpload) GetFileType(filePath string) (string, *systems.ErrorData) {
	buffer, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	filetype, err := filetype.Match(buffer)
	if err != nil {
		return "", ErrorMesg.InternalServerError(err.Error(), systems.CannotDetectFileType)
	}

	return filetype.Extension, nil
}

// SetCredential function used to set Amazon Credential
func (asu *AmazonS3ServiceUpload) SetCredential() (*credentials.Credentials, *systems.ErrorData) {
	creds := credentials.NewStaticCredentials(asu.AccessKey, asu.SecretKey, "")

	_, err1 := creds.Get()
	if err1 != nil {
		return nil, ErrorMesg.InternalServerError(err1.Error(), systems.ErrorAmazonService)
	}

	return creds, nil
}

// CreateSession function used to create new Amazon Session
func (asu *AmazonS3ServiceUpload) CreateSession() (*s3.S3, *systems.ErrorData) {
	credential, err := asu.SetCredential()
	if err != nil {
		return nil, err
	}
	// Initialize AWS Session with credentials
	awsSession := s3.New(session.New(), aws.NewConfig().WithRegion(asu.Region).WithCredentials(credential))

	return awsSession, nil
}
