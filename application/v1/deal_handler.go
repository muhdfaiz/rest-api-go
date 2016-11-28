package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DealHandler struct {
	DealService               DealServiceInterface
	DealTransformer           DealTransformerInterface
	ItemCategoryService       ItemCategoryServiceInterface
	DealCashbackService       DealCashbackServiceInterface
	UserRepository            UserRepositoryInterface
	ItemSubCategoryRepository ItemSubCategoryRepositoryInterface
}

// View function used to view deal details
func (dh *DealHandler) View(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	user := dh.UserRepository.GetByGUID(tokenData["user_guid"], "")

	// If user GUID not match user GUID inside the token return error message
	if user.GUID == "" {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals"))
		return
	}

	// Retrieve deal guid in url
	dealGUID := c.Param("deal_guid")

	// Retrieve query string for relations
	relations := c.Query("include")

	deal := dh.DealService.ViewDealDetails(dealGUID, relations)

	if deal.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal", "guid", dealGUID))
		return
	}

	// Check If deal quota still available for the user.
	total := dh.DealCashbackService.CountTotalNumberOfDealUserAddToList(user.GUID, deal.GUID)

	deal.CanAddTolist = 1

	if total >= deal.Perlimit {
		deal.CanAddTolist = 0
	}

	deal.RemainingAddToList = deal.Perlimit - total

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

	pageNumber := c.Query("page_number")
	pageLimit := c.Query("page_limit")
	name := c.Query("name")

	// Retrieve deals
	deals, totalDeal := dh.DealService.GetAvailableDealsForRegisteredUser(userGUID, name, latitude, longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, pageLimit)

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

	pageNumber := c.Query("page_number")
	pageLimit := c.Query("page_limit")

	// Retrieve query string parameter for relations
	relations := c.Query("include")

	// Retrieve deals
	deals, totalDeal := dh.DealService.GetAvailableDealsForGuestUser(latitude, longitude, pageNumber, pageLimit, relations)

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, pageLimit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewAndGroupByCategory function used to retrieve all deals group by category
func (dh *DealHandler) ViewAndGroupByCategory(c *gin.Context) {
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
	dealLimitPerCategory := c.Query("deal_limit_per_category")

	// Retrieve deals group by category
	dealCategories := dh.DealService.GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID, latitude, longitude, dealLimitPerCategory, "")

	c.JSON(http.StatusOK, gin.H{"data": dealCategories})
}

// ViewByCategoryAndGroupBySubCategory function used to retrieve all deals by category and group by subcategory
func (dh *DealHandler) ViewByCategoryAndGroupBySubCategory(c *gin.Context) {
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
		"latitude":  "required,latitude",
		"longitude": "required,longitude",
	}

	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), queryStringValidationRules)

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user guid in url
	categoryGUID := c.Param("category_guid")

	// Retrieve query string parameters for latitude and longitude
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	dealLimitPerSubcategory := c.Query("deal_limit_per_subcategory")

	// Retrieve deals group by category
	dealsGroupBySubCategory := dh.DealService.GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID, categoryGUID, latitude, longitude, dealLimitPerSubcategory, "")

	c.JSON(http.StatusOK, gin.H{"data": dealsGroupBySubCategory})
}

// ViewByCategory function used to retrieve all deals for specific category
func (dh *DealHandler) ViewByCategory(c *gin.Context) {
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

	pageNumber := c.Query("page_number")
	pageLimit := c.Query("page_limit")

	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by category"))
		return
	}

	// Retrieve category guid in url
	categoryGUID := c.Param("category_guid")

	dealCategory := dh.ItemCategoryService.GetItemCategoryByGUID(categoryGUID)

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

// ViewAllBySubCategory function used to retrieve all deals group by subcategory
func (dh *DealHandler) ViewBySubCategory(c *gin.Context) {
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

	pageNumber := c.Query("page_number")
	pageLimit := c.Query("page_limit")

	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by subcategory"))
		return
	}

	// Retrieve sub category guid in url
	subCategoryGUID := c.Param("subcategory_guid")

	dealSubCategory := dh.ItemSubCategoryRepository.GetByGUID(subCategoryGUID)

	// If user GUID not match user GUID inside the token return error message
	if dealSubCategory.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal Subcategory", "guid", subCategoryGUID))
		return
	}

	// Retrieve deals group by category
	deals, totalDeal := dh.DealService.GetAvailableDealsForSubCategoryForRegisteredUser(userGUID, subCategoryGUID, latitude,
		longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(c.Request, deals, totalDeal, pageLimit)

	c.JSON(http.StatusOK, gin.H{"data": result})
}
