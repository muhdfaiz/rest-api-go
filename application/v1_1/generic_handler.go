package v1_1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GenericHandler will handle all request related to generic category.
type GenericHandler struct {
	GenericService     GenericServiceInterface
	GenericTransformer GenericTransformerInterface
}

// ViewAll function used to retrieve all generic category through Generic Service.
func (gh *GenericHandler) ViewAll(context *gin.Context) {
	validationRules := map[string]string{
		"page_number":    "numeric",
		"page_limit":     "numeric",
		"last_sync_date": "time",
	}

	error := Validation.Validate(context.Request.URL.Query(), validationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	relations := context.Query("include")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	lastSyncDate := context.Query("last_sync_date")

	if lastSyncDate != "" {
		generics, totalGeneric := gh.GenericService.GetLatestUpdate(lastSyncDate, pageNumber, pageLimit, relations)

		genericResponse := gh.GenericTransformer.transformCollection(context.Request, generics, totalGeneric, pageLimit)

		context.JSON(http.StatusOK, genericResponse)
		return
	}

	generics, totalGeneric := gh.GenericService.GetAllGeneric(pageNumber, pageLimit, relations)

	genericResponse := gh.GenericTransformer.transformCollection(context.Request, generics, totalGeneric, pageLimit)

	context.JSON(http.StatusOK, genericResponse)
}
