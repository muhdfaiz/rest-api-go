package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DealHandler struct {
	DealService     DealServiceInterface
	DealTransformer DealTransformerInterface
}

// View function used to retrieve deal details
func (dh *DealHandler) View(c *gin.Context) {
	// Retrieve deal guid in url
	dealGUID := c.Param("guid")

	// Retrieve query string for relations
	relations := c.Query("include")

	deal := dh.DealService.ViewDealDetails(dealGUID, relations)

	if deal.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal", "guid", dealGUID))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deal})
}

// ViewAll function used to retrieve all deals based on latitude and longitude
func (dh *DealHandler) ViewAll(c *gin.Context) {
	// Validation Rules for query string parameters
	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), queryStringValidationRules)

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve query string parameters for latitude and longitude
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	offset := c.DefaultQuery("page_number", "1")
	limit := c.DefaultQuery("page_limit", "-1")

	// Retrieve query string parameter for relations
	relations := c.Query("include")

	// Retrieve deal ID
	deals, totalDeal := dh.DealService.GetAllDealsBasedOnLatitudeAndLongitudeAndQuota(latitude, longitude, offset, limit, relations)

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, limit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}
