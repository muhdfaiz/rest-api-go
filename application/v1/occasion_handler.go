package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OccasionHandler struct {
	OccasionRepository  OccasionRepositoryInterface
	OccasionTransformer OccasionTransformerInterface
}

// Index function used to retrieve shopping list occasions
func (oh *OccasionHandler) Index(c *gin.Context) {
	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), map[string]string{"last_sync_date": "time"})

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve filter query string in request
	lastSyncDate := c.DefaultQuery("last_sync_date", "")

	if lastSyncDate != "" {
		occasions, totalOccasion := oh.OccasionRepository.GetLatestUpdate(lastSyncDate)
		result := oh.OccasionTransformer.TransformCollection(occasions, totalOccasion)

		c.JSON(http.StatusOK, result)
		return
	}

	occasions, totalOccasion := oh.OccasionRepository.GetAll()
	result := oh.OccasionTransformer.TransformCollection(occasions, totalOccasion)

	c.JSON(http.StatusOK, result)
}
