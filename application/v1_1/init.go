package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

var (
	Database          = &systems.Database{}
	Error             = &systems.Error{}
	Helper            = &systems.Helpers{}
	Binding           = &systems.Binding{}
	FileValidation    = &filesystem.FileValidation{}
	FileSystem        = &filesystem.FileSystem{}
	PaginationReponse = &systems.PaginationResponse{}
	Validation        = &systems.Validation{}
	JWT               = &systems.Jwt{}
)
