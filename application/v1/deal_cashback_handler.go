package v1

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

// DealCashbackHandler will handle all requests related to Deal Cashback resource.
type DealCashbackHandler struct {
	ShoppingListRepository  ShoppingListRepositoryInterface
	DealCashbackService     DealCashbackServiceInterface
	DealCashbackTransformer DealCashbackTransformerInterface
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

	error := dch.DealCashbackService.CreateDealCashbackAndShoppingListItem(userGUID, dealCashbackData)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	result := make(map[string]string)
	result["message"] = "Successfully add deal guid " + dealCashbackData.DealGUID + " to list."

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// ViewByShoppingList function used to retrieve deal cashback by Shopping List GUID.
func (dch *DealCashbackHandler) ViewByShoppingList(context *gin.Context) {
	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view your deal cashbacks"))
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

	userDealCashbacks, totalUserDealCashback := dch.DealCashbackService.GetUserDealCashbackForUserShoppingList(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, relations)

	dealCashbackResponse := dch.DealCashbackTransformer.transformCollection(context.Request, userDealCashbacks, totalUserDealCashback, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": dealCashbackResponse})

}
