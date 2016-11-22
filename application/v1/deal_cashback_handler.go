package v1

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type DealCashbackHandler struct {
	ShoppingListRepository  ShoppingListRepositoryInterface
	DealCashbackService     DealCashbackServiceInterface
	DealCashbackTransformer DealCashbackTransformerInterface
}

// Create function used to add deal to the cashback and create shopping list item
func (dch *DealCashbackHandler) Create(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Bind request data based on header content type
	dealCashbackData := CreateDealCashback{}

	if err := Binding.Bind(&dealCashbackData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("add deal to list"))
		return
	}

	// Retrieve shopping list by guid and user guid
	shoppingList := dch.ShoppingListRepository.GetByGUIDAndUserGUID(dealCashbackData.ShoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", dealCashbackData.ShoppingListGUID))
		return
	}

	// Create deal cashback and shopping list item
	err := dch.DealCashbackService.CreateDealCashbackAndShoppingListItem(userGUID, dealCashbackData)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully add deal guid " + dealCashbackData.DealGUID + " to list."

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (dch *DealCashbackHandler) ViewByShoppingList(c *gin.Context) {
	// Validation Rules for query string parameters
	queryStringValidationRules := map[string]string{
		"page_number": "numeric",
		"page_limit":  "numeric",
	}

	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), queryStringValidationRules)

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view your deal cashbacks"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := dch.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	pageNumber := c.Query("page_number")
	pageLimit := c.Query("page_limit")

	// Retrieve query string for relations
	relations := c.Query("include")

	transactionStatus := c.Query("transaction_status")

	userDealCashbacks, totalUserDealCashback := dch.DealCashbackService.GetUserDealCashbackForUserShoppingList(userGUID, shoppingListGUID,
		transactionStatus, pageNumber, pageLimit, relations)

	dealCashbackResponse := dch.DealCashbackTransformer.transformCollection(c.Request, userDealCashbacks, totalUserDealCashback, pageLimit)

	c.JSON(http.StatusOK, gin.H{"data": dealCashbackResponse})

}
