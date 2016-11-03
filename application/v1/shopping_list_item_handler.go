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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve query string for relations
	relations := c.Query("include")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, relations)

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		DB.Rollback().Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	DB.Close()
	c.JSON(http.StatusOK, gin.H{"data": shoppingListItem})

}

// ViewAll function used to retrieve all user shopping list items
func (slih *ShoppingListItemHandler) ViewAll(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve added_to_cart query string param in url
	addedToCartBool, err := strconv.ParseBool(c.Query("added_to_cart"))

	if err != nil {
		// Retrieve shopping list item by guid
		userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItem(userGUID, shoppingListGUID, relations)
		DB.Close()
		c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
		return
	}

	if addedToCartBool == true {
		userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItemAddedToCart(userGUID, shoppingListGUID, relations)
		DB.Close()
		c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
		return
	}

	userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItemNotAddedToCart(userGUID, shoppingListGUID, relations)
	DB.Close()
	c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
	return

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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		DB.Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
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
	err := slih.ShoppingListItemFactory.UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID, shoppingListGUID, shoppingListItemGUID, structs.Map(updateShoppingListItemData))

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

// UpdateAll function used to update shopping list item
func (slih *ShoppingListItemHandler) UpdateAll(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
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
	err := slih.ShoppingListItemFactory.UpdateByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, structs.Map(updateShoppingListItemData))

	// If update shopping list item error, return error message
	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	DB.Commit()

	// Retrieve updated shopping list items
	result := slih.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	DB.Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// DeleteAll function used to delete all shopping list items including relation like image
// or to to delete all shopping list items from cart
func (slih *ShoppingListItemHandler) DeleteAll(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Soft delete Shopping List Item incuding relationship
	err1 := slih.ShoppingListItemFactory.DeleteByUserGUID(userGUID)

	if err1 != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err1)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted all shopping list item for user guid " + userGUID

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (slih *ShoppingListItemHandler) Delete(c *gin.Context) {
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		DB.Close()
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	// Soft delete Shopping List Item incuding relationship
	err := slih.ShoppingListItemFactory.DeleteByGUID(shoppingListItem.GUID)

	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted shopping list item with guid " + shoppingListItem.GUID

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": result})
}