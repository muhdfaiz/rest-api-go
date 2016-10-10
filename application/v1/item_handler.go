package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type ItemHandler struct {
	ItemRepository ItemRepositoryInterface
}

// Index function used to retrieve items
func (ih *ItemHandler) Index(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB)

	// Retrieve filter query string in request
	lastSyncDate := c.DefaultQuery("last_sync_date", "")

	if lastSyncDate != "" {
		items := ih.ItemRepository.GetLatestUpdate(lastSyncDate)
		DB.Close()
		c.JSON(http.StatusOK, gin.H{"last_update": items[len(items)-1].UpdatedAt, "data": items})
		return
	}

	items := ih.ItemRepository.GetAll()

	if len(items) == 0 {
		DB.Close()
		c.JSON(http.StatusOK, gin.H{"data": items})
		return
	}

	DB.Close()
	c.JSON(http.StatusOK, gin.H{"last_update": items[len(items)-1].UpdatedAt, "data": items})
}
