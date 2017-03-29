package v11

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GrocerHandler will handle all request related to Grocer resource.
type GrocerHandler struct {
	GrocerService GrocerServiceInterface
}

func (gh *GrocerHandler) GetAllGrocersThatContainDeals(context *gin.Context) {
	queryStringValidationRules := map[string]string{
		"latitude":  "required,latitude",
		"longitude": "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")

	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by subcategory"))
		return
	}

	grocers := gh.GrocerService.GetAllGrocersIncludingDeals(userGUID, latitude, longitude)

	context.JSON(http.StatusOK, gin.H{"data": grocers})
}
