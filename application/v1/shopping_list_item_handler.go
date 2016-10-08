package v1

import (
	"net/http"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ShoppingListItemHandler will handle all task related to resource shopping list item
type ShoppingListItemHandler struct {
	UserRepository               UserRepositoryInterface
	ShoppingListRepository       ShoppingListRepositoryInterface
	ShoppingListItemRepository   ShoppingListItemRepositoryInterface
	ShoppingListItemFactory      ShoppingListItemFactoryInterface
	ShoppingListItemImageFactory ShoppingListItemImageFactoryInterface
}

// View function used to retrieve shopping list items from database
func (slih *ShoppingListItemHandler) View(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB)
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		DB.Rollback().Close()
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, relations)

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	DB.Close()
	c.JSON(http.StatusOK, gin.H{"data": shoppingListItem})

}

// Create function used to create shopping list item
func (slih *ShoppingListItemHandler) Create(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		DB.Close()
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		DB.Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	createShoppingListItemData := CreateShoppingListItem{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&createShoppingListItemData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	createShoppingListItemData.ShoppingListGUID = shoppingListGUID
	createShoppingListItemData.UserGUID = userGUID

	// If shopping list item name already exist in the shopping list return error message
	shoppingListItem := slih.ShoppingListItemRepository.GetByName(createShoppingListItemData.Name, "")

	// If shopping list GUID empty return error message
	if shoppingListItem.Name != "" {
		DB.Close()
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("Shopping List Item", shoppingListItem.Name, "name"))
		return
	}

	// Create Shopping List item
	result, err := slih.ShoppingListItemFactory.Create(createShoppingListItemData)

	// Output error if failed to create new device
	if err != nil {
		DB.Rollback().Close()
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Update function used to update shopping list item
func (slih *ShoppingListItemHandler) Update(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		DB.Close()
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		DB.Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		DB.Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	updateShoppingListItemData := UpdateShoppingListItem{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&updateShoppingListItemData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Update Shopping List Item
	// structs.Map(updateShoppingListItemData)
	err := slih.ShoppingListItemFactory.Update(userGUID, shoppingListGUID, shoppingListItemGUID, structs.Map(updateShoppingListItemData))

	// If update shopping list item error, return error message
	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	DB.Commit()

	// Retrieve updated shopping list item
	result := slih.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	DB.Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}
