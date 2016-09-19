package v1

import (
	"bitbucket.org/shoppermate/services/facebook"
	"bitbucket.org/shoppermate/services/filesystem"
	"bitbucket.org/shoppermate/systems"
)

var (
	Database        = &systems.Database{}
	ErrorMesg       = &systems.Error{}
	Helper          = &systems.Helpers{}
	Config          = &systems.Configs{}
	Binding         = &systems.Binding{}
	Transformer     = &systems.Transformer{}
	FacebookService = &facebook.FacebookService{
		AppID:     Config.Get("app.yaml", "facebook_app_id", ""),
		AppSecret: Config.Get("app.yaml", "facebook_app_secret", ""),
	}
	FileValidation = &filesystem.FileValidation{}
	FileSystem     = &filesystem.FileSystem{}
)
