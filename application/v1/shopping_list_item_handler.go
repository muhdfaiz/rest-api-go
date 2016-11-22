package v1

import (
	"net/http"
	"strconv"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
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
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
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
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shoppingListItem})

}

// ViewAll function used to retrieve all user shopping list items
func (slih *ShoppingListItemHandler) ViewAll(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), map[string]string{"added_to_cart": "numeric", "latitude": "latitude", "longitude": "longitude"})

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list item"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guid
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve query string for latitude and longitude
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	// Retrieve added_to_cart query string param in url
	addedToCartBool, err1 := strconv.ParseBool(c.Query("added_to_cart"))

	if err1 != nil {
		// Retrieve shopping list item by guid
		userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItem(userGUID, shoppingListGUID, relations, latitude, longitude)

		c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
		return
	}

	if addedToCartBool == true {
		userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItemAddedToCart(userGUID, shoppingListGUID, relations)

		c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
		return
	}

	userShoppingListItems := slih.ShoppingListItemRepository.GetUserShoppingListItemNotAddedToCart(userGUID, shoppingListGUID, relations, latitude, longitude)

	c.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
	return

}

// Create function used to create shopping list item
func (slih *ShoppingListItemHandler) Create(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	createShoppingListItemData := CreateShoppingListItem{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&createShoppingListItemData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	createShoppingListItemData.ShoppingListGUID = shoppingListGUID
	createShoppingListItemData.UserGUID = userGUID

	// Create Shopping List item
	result, err := slih.ShoppingListItemFactory.Create(createShoppingListItemData)

	// Output error if failed to create new device
	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// Update function used to update shopping list item
func (slih *ShoppingListItemHandler) Update(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	updateShoppingListItemData := UpdateShoppingListItem{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&updateShoppingListItemData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Update Shopping List Item
	err := slih.ShoppingListItemFactory.UpdateByUserGUIDShoppingListGUIDAndShoppingListItemGUID(userGUID, shoppingListGUID, shoppingListItemGUID, structs.Map(updateShoppingListItemData))

	// If update shopping list item error, return error message
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve updated shopping list item
	result := slih.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// UpdateAll function used to update shopping list item
func (slih *ShoppingListItemHandler) UpdateAll(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	updateShoppingListItemData := UpdateShoppingListItem{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&updateShoppingListItemData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Update Shopping List Item
	err := slih.ShoppingListItemFactory.UpdateByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, structs.Map(updateShoppingListItemData))

	// If update shopping list item error, return error message
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve updated shopping list items
	result := slih.ShoppingListItemRepository.GetByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// DeleteAll function used to delete all shopping list items including relation like image
// or to to delete all shopping list items from cart
func (slih *ShoppingListItemHandler) DeleteAll(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), map[string]string{"added_to_cart": "numeric,len=1"})

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve only_added_to_cart query string in request
	addedToCart := c.Query("added_to_cart")

	if addedToCart == "1" {
		// Soft delete Shopping List Item incuding relationship
		err1 := slih.ShoppingListItemFactory.DeleteItemsHasBeenAddedToCartByUserGUID(userGUID)

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, err1)
			return
		}

		// Response data
		result := make(map[string]string)
		result["message"] = "Successfully deleted all shopping list items those has been added to cart for user guid " + userGUID

		c.JSON(http.StatusOK, gin.H{"data": result})
		return
	}

	if addedToCart == "0" {
		// Soft delete Shopping List Item incuding relationship
		err1 := slih.ShoppingListItemFactory.DeleteItemsHasNotBeenAddedToCartByUserGUID(userGUID)

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, err1)
			return
		}

		// Response data
		result := make(map[string]string)
		result["message"] = "Successfully deleted all shopping list items those hasn't been added to cart for user guid " + userGUID

		c.JSON(http.StatusOK, gin.H{"data": result})
		return
	}

	if addedToCart == "" {
		// Soft delete Shopping List Item incuding relationship
		err1 := slih.ShoppingListItemFactory.DeleteByUserGUID(userGUID)

		if err1 != nil {
			c.JSON(http.StatusInternalServerError, err1)
			return
		}
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted all shopping list item for user guid " + userGUID

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (slih *ShoppingListItemHandler) Delete(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	// Retrieve shopping list guid in url
	shoppingListGUID := c.Param("shopping_list_guid")

	// Retrieve shopping list by guid and user guif
	shoppingList := slih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := slih.ShoppingListItemRepository.GetByShoppingListGUIDAndGUID(shoppingListItemGUID, shoppingListGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	// Soft delete Shopping List Item incuding relationship
	err := slih.ShoppingListItemFactory.DeleteByGUID(shoppingListItem.GUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted shopping list item with guid " + shoppingListItem.GUID

	c.JSON(http.StatusOK, gin.H{"data": result})
}
