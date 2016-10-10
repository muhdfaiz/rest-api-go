package filesystem

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type AmazonS3Upload struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
}

// Upload function used to store file in local storage
func (asu *AmazonS3Upload) Upload(file multipart.File, localUploadPath string, amazonS3UploadPath string) (map[string]string, *systems.ErrorData) {

	localUploadService := LocalUpload{}
	uploadedFile, err := localUploadService.Upload(file, localUploadPath)
	if err != nil {
		return nil, err
	}

	// Read file
	localFile, err1 := os.Open(localUploadPath + uploadedFile["name"])

	if err1 != nil {
		return nil, Error.InternalServerError(err1.Error(), systems.CannotReadFile)
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
	path := amazonS3UploadPath + uploadedFile["name"]
	params := &s3.PutObjectInput{
		Bucket:        aws.String(asu.BucketName),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(uploadedFile["type"]),
	}

	awsSession, err := asu.createSession()
	if err != nil {
		return nil, err
	}

	_, err2 := awsSession.PutObject(params)
	if err2 != nil {
		return nil, Error.InternalServerError(err2.Error(), systems.ErrorAmazonService)
	}

	// Remove Local File
	err2 = os.Remove(uploadedFile["path"])
	if err2 != nil {
		return nil, Error.InternalServerError(err2.Error(), systems.CannotDeleteFile)
	}

	uploadedFile["path"] = fmt.Sprintf("https://s3-%s.amazonaws.com/%s%s%s", asu.Region,
		asu.BucketName, amazonS3UploadPath, uploadedFile["name"])

	return uploadedFile, nil
}

// Delete function used to delete multiple file in Amazon S3
func (asu *AmazonS3Upload) Delete(files []string) *systems.ErrorData {
	amazonS3Objects := make([]*s3.ObjectIdentifier, len(files))

	for key, file := range files {
		amazonS3Objects[key] = &s3.ObjectIdentifier{Key: aws.String(file)}
	}

	awsSession, err := asu.createSession()

	if err != nil {
		return err
	}

	params := &s3.DeleteObjectsInput{
		Bucket: aws.String(asu.BucketName),
		Delete: &s3.Delete{
			Objects: amazonS3Objects,
		},
	}

	_, err1 := awsSession.DeleteObjects(params)

	if err1 != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return Error.InternalServerError(err1.Error(), systems.FailedToDeleteAmazonS3File)
	}

	return nil
}

// SetCredential function used to set Amazon Credential
func (asu *AmazonS3Upload) setCredential() (*credentials.Credentials, *systems.ErrorData) {
	creds := credentials.NewStaticCredentials(asu.AccessKey, asu.SecretKey, "")

	_, err1 := creds.Get()
	if err1 != nil {
		return nil, Error.InternalServerError(err1.Error(), systems.ErrorAmazonService)
	}

	return creds, nil
}

// CreateSession function used to create new Amazon Session
func (asu *AmazonS3Upload) createSession() (*s3.S3, *systems.ErrorData) {
	credential, err := asu.setCredential()
	if err != nil {
		return nil, err
	}
	// Initialize AWS Session with credentials
	awsSession := s3.New(session.New(), aws.NewConfig().WithRegion(asu.Region).WithCredentials(credential))

	return awsSession, nil
}
