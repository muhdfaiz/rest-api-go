package filesystem

import (
	"io"
	"mime/multipart"
	"os"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type LocalUpload struct {
}

type LocalUploadConfigInterface interface {
	GetLocalUploadPath() string
}

// Upload function used to store file in local storage
func (lu *LocalUpload) Upload(file multipart.File, uploadPath string) (map[string]string, *systems.ErrorData) {
	// Generate Unique File Name
	fileSystem := &FileSystem{}

	fileType, err := fileSystem.DetectFileType(file)
	if err != nil {
		return nil, err
	}
	fileName := fileSystem.GenerateUniqueFileName() + "." + fileType

	// Create new file in local storage
	localFile, err1 := os.Create(uploadPath + fileName)

	if err1 != nil {
		return nil, Error.InternalServerError(err1.Error(), systems.CannotReadFile)
	}
	defer localFile.Close()

	// Copy file from request to created file in local storage
	_, err1 = io.Copy(localFile, file)
	if err1 != nil {
		return nil, Error.InternalServerError(err1.Error(), systems.CannotCopyFile)
	}

	result := make(map[string]string)
	result["name"] = fileName
	result["path"] = uploadPath + fileName
	result["type"] = fileType

	return result, nil
}
