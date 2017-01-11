package v1_1

import (
	"net/http"
	"strconv"

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
func (dh *DealHandler) View(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)
	dealGUID := context.Param("deal_guid")

	relations := context.Query("include")

	deal := dh.DealService.ViewDealDetails(dealGUID, relations)

	if deal.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal", "guid", dealGUID))
		return
	}

	total := dh.DealCashbackService.CountTotalNumberOfDealUserAddedToList(tokenData["user_guid"], deal.GUID)

	deal.CanAddTolist = 1

	if total >= deal.Perlimit {
		deal.CanAddTolist = 0
	}

	deal.RemainingAddToList = deal.Perlimit - total

	context.JSON(http.StatusOK, gin.H{"data": deal})
}

// ViewAllForRegisteredUser function used to retrieve all deals based on latitude and longitude
func (dh *DealHandler) ViewAllForRegisteredUser(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals"))
		return
	}

	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	name := context.Query("name")

	// Retrieve deals
	deals, totalDeal := dh.DealService.GetAvailableDealsForRegisteredUser(userGUID, name, latitude, longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(context.Request, deals, totalDeal, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewAllForGuestUser function used to retrieve all deals based on latitude and longitude for guest user
// The different between ViewAllForGuestUser and ViewAllForRegisteredUser is ViewAllForRegisteredUser function
// will validate how many time user already add same deal to list. API need to use User GUID to validate that information.
func (dh *DealHandler) ViewAllForGuestUser(context *gin.Context) {
	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	relations := context.Query("include")

	deals, totalDeal := dh.DealService.GetAvailableDealsForGuestUser(latitude, longitude, pageNumber, pageLimit, relations)

	result := dh.DealTransformer.transformCollection(context.Request, deals, totalDeal, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewAndGroupByCategory function used to retrieve all deals group by category
func (dh *DealHandler) ViewAndGroupByCategory(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals by category"))
		return
	}

	queryStringValidationRules := map[string]string{
		"latitude":                "required,latitude",
		"longitude":               "required,longitude",
		"deal_limit_per_category": "numeric",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	dealLimitPerCategory := context.Query("deal_limit_per_category")

	dealCategories := dh.DealService.GetAvailableDealsGroupByCategoryForRegisteredUser(userGUID, latitude, longitude, dealLimitPerCategory, "")

	context.JSON(http.StatusOK, gin.H{"data": dealCategories})
}

// ViewByCategoryAndGroupBySubCategory function used to retrieve all deals by category and group by subcategory
func (dh *DealHandler) ViewByCategoryAndGroupBySubCategory(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := context.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deals by category"))
		return
	}

	// Validation Rules for query string parameters
	queryStringValidationRules := map[string]string{
		"latitude":  "required,latitude",
		"longitude": "required,longitude",
	}

	// Validate query string
	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	// If validation error return error message
	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	// Retrieve user guid in url
	categoryGUID := context.Param("category_guid")

	// Retrieve query string parameters for latitude and longitude
	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	dealLimitPerSubcategory := context.Query("deal_limit_per_subcategory")

	// Retrieve deals group by category
	dealsGroupBySubCategory := dh.DealService.GetAvailableDealsByCategoryGroupBySubCategoryForRegisteredUser(userGUID, categoryGUID, latitude, longitude, dealLimitPerSubcategory, "")

	context.JSON(http.StatusOK, gin.H{"data": dealsGroupBySubCategory})
}

// ViewByCategory function used to retrieve all deals for specific category
func (dh *DealHandler) ViewByCategory(context *gin.Context) {
	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")

	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by category"))
		return
	}

	categoryGUID := context.Param("category_guid")

	dealCategory := dh.ItemCategoryService.GetItemCategoryByGUID(categoryGUID)

	if dealCategory.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal Category", "guid", categoryGUID))
		return
	}

	deals, totalDeal := dh.DealService.GetAvailableDealsByCategoryForRegisteredUser(userGUID, dealCategory.Name, latitude,
		longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(context.Request, deals, totalDeal, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewBySubCategory function used to retrieve all deals group by subcategory
func (dh *DealHandler) ViewBySubCategory(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by subcategory"))
		return
	}

	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")

	subCategoryGUID := context.Param("subcategory_guid")

	dealSubCategory := dh.ItemSubCategoryRepository.GetByGUID(subCategoryGUID)

	if dealSubCategory.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Deal Subcategory", "guid", subCategoryGUID))
		return
	}

	deals, totalDeal := dh.DealService.GetAvailableDealsForSubCategoryForRegisteredUser(userGUID, subCategoryGUID, latitude,
		longitude, pageNumber, pageLimit, "")

	result := dh.DealTransformer.transformCollection(context.Request, deals, totalDeal, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewByGrocerAndCategory function used to retrieve valid deals group by subcategory
func (dh *DealHandler) ViewByGrocerAndCategory(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)
	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all deals by subcategory"))
		return
	}

	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
		"latitude":    "required,latitude",
		"longitude":   "required,longitude",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")

	grocerGUID := context.Param("grocer_guid")
	categoryGUID := context.Param("category_guid")

	dealResponse, error := dh.DealService.GetAvailableDealsForGrocerByCategory(context.Request, userGUID,
		grocerGUID, categoryGUID, latitude, longitude, pageNumber, pageLimit, "")

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, dealResponse)
}
