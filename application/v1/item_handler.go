package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ItemHandler struct {
	ItemRepository  ItemRepositoryInterface
	ItemTransformer ItemTransformerInterface
}

// Index function used to retrieve items list from database
func (ih *ItemHandler) Index(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB)

	// Validate query string
	err := Validation.Validate(c, map[string]string{"last_sync_date": "time", "page_number": "numeric", "page_limit": "numeric"})

	// If validation error return error message
	if err != nil {
		DB.Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve query string in request
	lastSyncDate := c.DefaultQuery("last_sync_date", "")
	offset := c.DefaultQuery("page_number", "1")
	limit := c.DefaultQuery("page_limit", "-1")

	if lastSyncDate != "" {
		items, totalItems := ih.ItemRepository.GetLatestUpdate(lastSyncDate, offset, limit)
		result := ih.ItemTransformer.transformCollection(c.Request, items, totalItems, limit)

		DB.Close()
		c.JSON(http.StatusOK, result)
		return
	}

	items, totalItems := ih.ItemRepository.GetAll(offset, limit)

	result := ih.ItemTransformer.transformCollection(c.Request, items, totalItems, limit)
	DB.Close()

	c.JSON(http.StatusOK, result)
}
