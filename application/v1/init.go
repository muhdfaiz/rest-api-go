package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

var (
	Database        = &systems.Database{}
	Error           = &systems.Error{}
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
