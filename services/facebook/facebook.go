package facebook

import (
	"fmt"
	"os"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	fb "github.com/huandu/facebook"
)

var (
	Error  = &systems.Error{}
	Helper = &systems.Helpers{}
)

type FacebookServiceInterface interface {
	IDIsValid(facebookID string, debug int) bool
	GetAccessToken(debug int) string
}

// FacebookService type
type FacebookService struct{}

// IDIsValid function used to  verify Facebook ID valid or not
// Return true if valid otherwise return false
func (fs *FacebookService) IDIsValid(facebookID string, debug int) bool {
	fbAccessToken := fs.GetAccessToken(debug)

	result, _ := fb.Get(fmt.Sprintf("/%s", facebookID), fb.Params{
		"access_token": fbAccessToken,
	})

	if _, ok := result["id"]; ok {
		return true
	}

	return false
}

// GetAccessToken function used to retrieve access token
func (fs *FacebookService) GetAccessToken(debug int) string {
	appID := os.Getenv("FACEBOOK_APP_ID")

	appSecret := os.Getenv("FACEBOOK_APP_SECRET")

	if debug == 1 {
		appID = os.Getenv("DEBUG_FACEBOOK_APP_ID")

		appSecret = os.Getenv("DEBUG_FACEBOOK_APP_SECRET")
	}

	fbApp := fb.New(appID, appSecret)

	return fbApp.AppAccessToken()
}
