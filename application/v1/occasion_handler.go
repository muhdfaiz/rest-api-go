package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type OccasionHandler struct {
	OccasionRepository OccasionRepositoryInterface
}

// Index function used to retrieve shopping list occasions
func (oh *OccasionHandler) Index(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve filter query string in request
	lastSyncDate := c.DefaultQuery("last_sync_date", "")

	if lastSyncDate != "" {
		occasions := oh.OccasionRepository.GetLatestUpdate(lastSyncDate)
		DB.Commit().Close()
		c.JSON(http.StatusOK, gin.H{"last_update": occasions[len(occasions)-1].UpdatedAt, "data": occasions})
		return
	}

	occasions := oh.OccasionRepository.GetAll()

	if len(occasions) == 0 {
		DB.Commit().Close()
		c.JSON(http.StatusOK, gin.H{"data": occasions})
		return
	}

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"last_update": occasions[len(occasions)-1].UpdatedAt, "data": occasions})
}
