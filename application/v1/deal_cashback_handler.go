package v1

import (
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DealCashbackHandler struct {
	ShoppingListRepository ShoppingListRepositoryInterface
	DealCashbackService    DealCashbackServiceInterface
}

// Create function used to add deal to the cashback and create shopping list item
func (dch *DealCashbackHandler) Create(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()
	tokenData := c.MustGet("Token").(map[string]string)

	// Bind request data based on header content type
	dealCashbackData := CreateDealCashback{}

	if err := Binding.Bind(&dealCashbackData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		DB.Close()
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("add deal to list"))
		return
	}

	// Retrieve shopping list by guid and user guid
	shoppingList := dch.ShoppingListRepository.GetByGUIDAndUserGUID(dealCashbackData.ShoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		DB.Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", dealCashbackData.ShoppingListGUID))
		return
	}

	// Create deal cashback and shopping list item
	err := dch.DealCashbackService.CreateDealCashbackAndShoppingListItem(userGUID, dealCashbackData)

	if err != nil {
		DB.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully add deal guid " + dealCashbackData.DealGUID + "to list."

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
