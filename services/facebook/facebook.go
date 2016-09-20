package facebook

import (
	"fmt"

	"bitbucket.org/shoppermate-api/systems"
	fb "github.com/huandu/facebook"
)

var (
	Config = &systems.Configs{}
	Error  = &systems.Error{}
	Helper = &systems.Helpers{}
)

// FacebookService type
type FacebookService struct {
	AppID     string
	AppSecret string
}

// IDIsValid function used to  verify Facebook ID valid or not
// Return true if valid otherwise return false
func (fs *FacebookService) IDIsValid(facebookID string) bool {
	// If facebook App ID or Secret empty, retrieve App ID & Secret from config
	fbAccessToken := fs.GetAccessToken()

	result, _ := fb.Get(fmt.Sprintf("/%s", facebookID), fb.Params{
		"access_token": fbAccessToken,
	})

	if _, ok := result["id"]; ok {
		return true
	}

	return false
}

// GetAccessToken function used to retrieve access token
func (fs *FacebookService) GetAccessToken() string {
	fbApp := fb.New(fs.AppID, fs.AppSecret)

	return fbApp.AppAccessToken()
}
