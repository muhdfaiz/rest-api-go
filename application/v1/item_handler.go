package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	ItemRepository  ItemRepositoryInterface
	ItemTransformer ItemTransformerInterface
}

// Index function used to retrieve items list from database
func (ih *ItemHandler) Index(c *gin.Context) {
	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), map[string]string{"last_sync_date": "time", "page_number": "numeric", "page_limit": "numeric"})

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve query string in request
	lastSyncDate := c.DefaultQuery("last_sync_date", "")
	offset := c.DefaultQuery("page_number", "1")
	limit := c.DefaultQuery("page_limit", "-1")
	relations := c.DefaultQuery("include", "")

	if lastSyncDate != "" {
		items, totalItems := ih.ItemRepository.GetLatestUpdate(lastSyncDate, offset, limit, relations)
		result := ih.ItemTransformer.transformCollection(c.Request, items, totalItems, limit)

		c.JSON(http.StatusOK, result)
		return
	}

	items, totalItems := ih.ItemRepository.GetAll(offset, limit, relations)

	result := ih.ItemTransformer.transformCollection(c.Request, items, totalItems, limit)

	c.JSON(http.StatusOK, result)
}

// GetCategories function used to retrieve unique item categories list from database
func (ih *ItemHandler) GetCategories(c *gin.Context) {
	itemCategories, totalItemCategories := ih.ItemRepository.GetUniqueCategories("")

	c.JSON(http.StatusOK, gin.H{"total_data": totalItemCategories, "data": itemCategories})
}
