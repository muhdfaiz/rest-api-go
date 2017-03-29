package v11

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefaultShoppingListHandler will handle all request related to Default Shopping List resource.
type DefaultShoppingListHandler struct {
	DefaultShoppingListService DefaultShoppingListServiceInterface
}

// ViewAll function used to retrieve all default shopping lists.
func (dslh *DefaultShoppingListHandler) ViewAll(context *gin.Context) {
	error := Validation.Validate(context.Request.URL.Query(), map[string]string{"latitude": "latitude", "longitude": "longitude"})

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")

	longitude := context.Query("longitude")

	defaultShoppingLists := dslh.DefaultShoppingListService.GetAllDefaultShoppingListsIncludingItemsAndDeals(latitude, longitude, "occasions,items.images")

	context.JSON(http.StatusOK, gin.H{"data": defaultShoppingLists})
}
