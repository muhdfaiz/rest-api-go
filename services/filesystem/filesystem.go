package filesystem

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"bitbucket.org/shoppermate/systems"
	filetype "gopkg.in/h2non/filetype.v0"
)

var (
	ErrorMesg = &systems.Error{}
)

// type FileSystemInterface interface {
// 	Upload(lci *LocalConfigInterface)
// 	Delete(path string)
// 	List(path string)
// 	CreateDir(dirName string, configs map[string]string)
// 	DeleteDir(dirname string)
// 	FormatOutput()
// }

type FileSystem struct{}

type FileValidation struct{}

func (fs FileSystem) Driver(driver string) interface{} {
	switch driver {
	case "local":
		return &LocalUpload{}
	case "amazonS3":
		return &AmazonS3ServiceUpload{}
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
		return "", ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	filetype, err := filetype.Match(bytes)
	if err != nil {
		return "", ErrorMesg.InternalServerError(err.Error(), systems.CannotDetectFileType)
	}

	return filetype.Extension, nil
}

// ValidateFileType function used to validate file type allow to upload
func (fv *FileValidation) ValidateFileType(allowFileTypes []string, file multipart.File) (string, *systems.ErrorData) {
	// Skip validation file type if FileTypeAllow empty
	if allowFileTypes == nil {
		return "", nil
	}

	// Detect filetype
	fileSystem := FileSystem{}
	fileType, err := fileSystem.DetectFileType(file)
	if err != nil {
		return "", err
	}

	for _, allowFileType := range allowFileTypes {
		if allowFileType == fileType {
			return fileType, nil
		}
	}

	return "", ErrorMesg.InvalidFileTypeError(strings.Join(allowFileTypes, ", "))
}

// ValidateFileSize function used to verify if file size want to upload is not bigger than system allowed
func (fv *FileValidation) ValidateFileSize(file multipart.File, maxFileSizeAllow int64, fileAttribute string) (int64, *systems.ErrorData) {
	fileSize, err := file.Seek(0, 2)
	if err != nil {
		return 0, ErrorMesg.InternalServerError(err.Error(), systems.CannotReadFile)
	}

	if maxFileSizeAllow == 0 {
		return fileSize, nil
	}

	if fileSize > maxFileSizeAllow {
		return fileSize, ErrorMesg.FileSizeExceededLimit(fileAttribute, strconv.FormatInt(maxFileSizeAllow/1000, 10))
	}
	file.Seek(0, 0)
	return fileSize, nil
}
