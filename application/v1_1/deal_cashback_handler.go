package v1

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DealCashbackHandler will handle all requests related to Deal Cashback resource.
type DealCashbackHandler struct {
	ShoppingListRepository  ShoppingListRepositoryInterface
	DealCashbackService     DealCashbackServiceInterface
	DealCashbackTransformer DealCashbackTransformerInterface
	DealService             DealServiceInterface
}

// Create function used to create new deal cashback and store in database and create shopping list item based on deal info.
func (dch *DealCashbackHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	dealCashbackData := CreateDealCashback{}

	if error := Binding.Bind(&dealCashbackData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("add deal to list"))
		return
	}

	shoppingList := dch.ShoppingListRepository.GetByGUIDAndUserGUID(dealCashbackData.ShoppingListGUID, userGUID, "")

	if shoppingList.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", dealCashbackData.ShoppingListGUID))
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	error := dch.DealCashbackService.CreateDealCashbackAndShoppingListItem(dbTransaction, userGUID, dealCashbackData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	result := make(map[string]string)
	result["message"] = "Successfully add deal guid " + dealCashbackData.DealGUID + " to list."

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewByShoppingList function used to retrieve deal cashback by Shopping List GUID.
func (dch *DealCashbackHandler) ViewByShoppingList(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deal cashbacks in shopping list"))
		return
	}

	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	shoppingList := dch.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	if shoppingList.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	relations := context.Query("include")

	transactionStatus := context.Query("transaction_status")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	userDealCashbacks, totalUserDealCashback, error := dch.DealCashbackService.GetUserDealCashbacksByShoppingList(dbTransaction, userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, relations)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dealCashbackResponse := dch.DealCashbackTransformer.transformCollection(context.Request, userDealCashbacks, totalUserDealCashback, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": dealCashbackResponse})
}

// ViewByUserAndDealGroupByShoppingList function used to retrieve all deal cashbacks by user GUID and deal GUID
// including shopping list and group the result by Shopping List.
func (dch *DealCashbackHandler) ViewByUserAndDealGroupByShoppingList(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deal cashbacks"))
		return
	}

	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	dealGUID := context.Param("deal_guid")

	_, error = dch.DealService.CheckDealExistOrNotByGUID(dealGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")

	dealCashbacks, totalDealCashbacks := dch.DealCashbackService.GetUserDealCashbacksByDealGUID(userGUID, dealGUID, pageNumber, pageLimit, "shoppinglists")

	dealCashbackResponse := dch.DealCashbackTransformer.transformCollection(context.Request, dealCashbacks, totalDealCashbacks, pageLimit)

	context.JSON(http.StatusOK, dealCashbackResponse)

}
