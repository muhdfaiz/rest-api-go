package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DealHandler struct {
	DealService         DealServiceInterface
	DealTransformer     DealTransformerInterface
	ItemCategoryService ItemCategoryServiceInterface
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

// ViewAllForRegisteredUser function used to retrieve all deals based on latitude and longitude
func (dh *DealHandler) ViewAllForRegisteredUser(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals"))
		return
	}

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

	// Retrieve deals
	deals, totalDeal := dh.DealService.GetAvailableDealsForRegisteredUser(userGUID, latitude, longitude, offset, limit, relations)

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, limit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewAllForGuestUser function used to retrieve all deals based on latitude and longitude for guest user
// The different between ViewAllForGuestUser and ViewAllForRegisteredUser is ViewAllForRegisteredUser function
// will validate how many time user already add same deal to list. API need to use User GUID to validate that information.
func (dh *DealHandler) ViewAllForGuestUser(c *gin.Context) {
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

	// Retrieve deals
	deals, totalDeal := dh.DealService.GetAvailableDealsForGuestUser(latitude, longitude, offset, limit, relations)

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, limit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewAllGroupByCategory function used to retrieve all deals group by category
func (dh *DealHandler) ViewAllGroupByCategory(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals by category"))
		return
	}

	// Validation Rules for query string parameters
	queryStringValidationRules := map[string]string{
		"latitude":                "required,latitude",
		"longitude":               "required,longitude",
		"deal_limit_per_category": "numeric",
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
	dealLimitPerCategory := c.DefaultQuery("deal_limit_per_category", "10000")

	// Retrieve deals group by category
	dealCategories := dh.DealService.GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID, latitude, longitude, "1", dealLimitPerCategory, "")

	c.JSON(http.StatusOK, gin.H{"data": dealCategories})
}

// ViewAllByCategory function used to retrieve all deals group by category
func (dh *DealHandler) ViewAllByCategory(c *gin.Context) {
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

	pageNumber := c.DefaultQuery("page_number", "1")
	pageLimit := c.DefaultQuery("page_limit", "-1")

	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals by category"))
		return
	}

	// Retrieve category guid in url
	categoryGUID := c.Param("category_guid")

	dealCategory := dh.ItemCategoryService.GetItemCategoByGUID(categoryGUID)

	// If user GUID not match user GUID inside the token return error message
	if dealCategory.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal Category", "guid", categoryGUID))
		return
	}

	// Retrieve deals group by category
	deals, totalDeal := dh.DealService.GetAvailableDealsByCategoryForRegisteredUser(userGUID, dealCategory.Name, latitude,
		longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, pageLimit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}
