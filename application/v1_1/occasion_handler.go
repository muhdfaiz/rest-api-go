package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OccasionHandler struct {
	OccasionService OccasionServiceInterface
}

// Index function used to retrieve shopping list occasions
func (oh *OccasionHandler) Index(context *gin.Context) {
	error := Validation.Validate(context.Request.URL.Query(), map[string]string{"last_sync_date": "time"})

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	lastSyncDate := context.Query("last_sync_date")

	if lastSyncDate != "" {
		occasions := oh.OccasionService.GetLatestOccasionAfterLastSyncDate(lastSyncDate)

		context.JSON(http.StatusOK, occasions)
		return
	}

	occasions := oh.OccasionService.GetAllOccasions()

	context.JSON(http.StatusOK, occasions)
}
