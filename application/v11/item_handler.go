package v11

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ItemHandler will handle all request related to Item resource.
type ItemHandler struct {
	ItemRepository  ItemRepositoryInterface
	ItemTransformer ItemTransformerInterface
}

// Index function used to retrieve items list from database
func (ih *ItemHandler) Index(context *gin.Context) {
	error := Validation.Validate(context.Request.URL.Query(), map[string]string{"last_sync_date": "time", "page_number": "numeric", "page_limit": "numeric"})

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	// Retrieve query string in request
	lastSyncDate := context.DefaultQuery("last_sync_date", "")
	offset := context.DefaultQuery("page_number", "1")
	limit := context.DefaultQuery("page_limit", "-1")
	relations := context.DefaultQuery("include", "")

	if lastSyncDate != "" {
		items, totalItems := ih.ItemRepository.GetLatestUpdate(lastSyncDate, offset, limit, relations)
		result := ih.ItemTransformer.transformCollection(context.Request, items, totalItems, limit)

		context.JSON(http.StatusOK, result)
		return
	}

	items, totalItems := ih.ItemRepository.GetAll(offset, limit, relations)

	result := ih.ItemTransformer.transformCollection(context.Request, items, totalItems, limit)

	context.JSON(http.StatusOK, result)
}

// GetCategories function used to retrieve unique item categories list from database
func (ih *ItemHandler) GetCategories(context *gin.Context) {
	itemCategories, totalItemCategories := ih.ItemRepository.GetUniqueCategories("")

	context.JSON(http.StatusOK, gin.H{"total_data": totalItemCategories, "data": itemCategories})
}
