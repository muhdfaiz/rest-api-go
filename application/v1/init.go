package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

var (
	Database          = &systems.Database{}
	Error             = &systems.Error{}
	Helper            = &systems.Helpers{}
	Config            = &systems.Configs{}
	Binding           = &systems.Binding{}
	FileValidation    = &filesystem.FileValidation{}
	FileSystem        = &filesystem.FileSystem{}
	PaginationReponse = &systems.PaginationResponse{}
	Validation        = &systems.Validation{}
)
