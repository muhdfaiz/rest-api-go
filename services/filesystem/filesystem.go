package filesystem

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	uuid "github.com/satori/go.uuid"
	filetype "gopkg.in/h2non/filetype.v0"
)

type FileSystem struct{}

type FileValidation struct{}

func (fs *FileSystem) Driver(driver string) interface{} {
	switch driver {
	case "local":
		return &LocalUpload{}
	case "amazonS3":
		return &AmazonS3Upload{}
	}
	return &LocalUpload{}
}

// GenerateUniqueFileName function used to generate unique file name using uuid v5 with timestamp
func (fs *FileSystem) GenerateUniqueFileName() string {
	namespace := uuid.NamespaceDNS
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return uuid.NewV3(namespace, timestamp).String()
}

// DetectFileType function will validate file type allow to upload or not
func (fs *FileSystem) DetectFileType(file multipart.File) (string, *systems.ErrorData) {
	bytes := make([]byte, 4)

	_, err := file.ReadAt(bytes, 0)
	if err != nil {
		return "", Error.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	filetype, err := filetype.Match(bytes)
	if err != nil {
		return "", Error.InternalServerError(err.Error(), systems.CannotDetectFileType)
	}

	return filetype.Extension, nil
}

// ValidateFileType function used to validate file type allow to upload
func (fv *FileValidation) ValidateFileType(allowFileTypes []string, file multipart.File) *systems.ErrorData {
	// Skip validation file type if FileTypeAllow empty
	if len(allowFileTypes) <= 0 {
		return nil
	}

	// Detect filetype
	fileSystem := FileSystem{}
	fileType, err := fileSystem.DetectFileType(file)
	if err != nil {
		return err
	}

	for _, allowFileType := range allowFileTypes {
		if allowFileType == fileType {
			return nil
		}
	}

	return Error.InvalidFileTypeError(strings.Join(allowFileTypes, ", "))
}

// ValidateFileSize function used to verify if file size want to upload is not bigger than system allowed
func (fv *FileValidation) ValidateFileSize(file multipart.File, maxFileSizeAllow int64, fileAttribute string) (int64, *systems.ErrorData) {
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return 0, Error.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	if maxFileSizeAllow == 0 {
		return fileSize, nil
	}

	if fileSize > maxFileSizeAllow {
		return fileSize, Error.FileSizeExceededLimit(fileAttribute, strconv.FormatInt(maxFileSizeAllow/1000, 10))
	}

	file.Seek(0, 0)
	return fileSize, nil
}
