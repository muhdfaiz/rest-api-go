package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShoppingListItemHandler will handle all task related to shopping list item resource.
type ShoppingListItemHandler struct {
	ShoppingListItemService      ShoppingListItemServiceInterface
	ShoppingListService          ShoppingListServiceInterface
	ShoppingListItemImageService ShoppingListItemImageServiceInterface
}

// View function used to retrieve user shopping list items from database.
func (slih *ShoppingListItemHandler) View(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	shoppingListItemGUID := context.Param("item_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view shopping list"))
		return
	}

	_, error := slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	shoppingListItem, error := slih.ShoppingListItemService.ViewUserShoppingListItem(userGUID, shoppingListGUID, shoppingListItemGUID, "occasions,items.images")

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": shoppingListItem})
}

// ViewAll function used to retrieve all user shopping list items.
func (slih *ShoppingListItemHandler) ViewAll(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list item"))
		return
	}

	error := Validation.Validate(context.Request.URL.Query(), map[string]string{"added_to_cart": "numeric", "latitude": "latitude", "longitude": "longitude"})

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	relations := context.Query("include")

	latitude := context.Query("latitude")

	longitude := context.Query("longitude")

	addedToCart := context.Query("added_to_cart")

	_, error = slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	userShoppingListItems, error := slih.ShoppingListItemService.ViewAllUserShoppingListItem(userGUID, shoppingListGUID, addedToCart, latitude, longitude, relations)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
}

// Create function used to create user shopping list item
func (slih *ShoppingListItemHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	shoppingListItemToCreate := CreateShoppingListItem{}

	if error := Binding.Bind(&shoppingListItemToCreate, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	_, error := slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	createdShoppingListItem, error := slih.ShoppingListItemService.CreateUserShoppingListItem(userGUID, shoppingListGUID, shoppingListItemToCreate)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": createdShoppingListItem})
}

// Update function used to update one of the user shopping list item
func (slih *ShoppingListItemHandler) Update(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	shoppingListItemToUpdate := UpdateShoppingListItem{}

	if error := Binding.Bind(&shoppingListItemToUpdate, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")
	shoppingListItemGUID := context.Param("item_guid")

	_, error := slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	updatedShoppingList, error := slih.ShoppingListItemService.UpdateUserShoppingListItem(userGUID, shoppingListGUID, shoppingListItemGUID, shoppingListItemToUpdate, "")

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": updatedShoppingList})
}

// UpdateAll function used to update shopping list item
func (slih *ShoppingListItemHandler) UpdateAll(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	updateShoppingListItemData := UpdateShoppingListItem{}

	if error := Binding.Bind(&updateShoppingListItemData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	_, error := slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	updatedShoppingListItems, error := slih.ShoppingListItemService.UpdateAllUserShoppingListItem(userGUID, shoppingListGUID, updateShoppingListItemData)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": updatedShoppingListItems})
}

// DeleteAll function used to delete all shopping list items including relation like image
// or to to delete all shopping list items from cart
func (slih *ShoppingListItemHandler) DeleteAll(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	error := Validation.Validate(context.Request.URL.Query(), map[string]string{"added_to_cart": "numeric,len=1"})

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	deleteItemInCart := context.Query("added_to_cart")

	_, error = slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	result, error := slih.ShoppingListItemService.DeleteUserShoppingListItem(userGUID, shoppingListGUID, deleteItemInCart)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	error = slih.ShoppingListItemImageService.DeleteImagesForShoppingList(shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// Delete function used to delete user shopping list item by shopping list item GUID, user GUID and shopping list GUID.
func (slih *ShoppingListItemHandler) Delete(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("delete shopping list"))
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")
	shoppingListItemGUID := context.Param("item_guid")

	result, error := slih.ShoppingListItemService.DeleteShoppingListItemInShoppingList(shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	error = slih.ShoppingListItemImageService.DeleteImagesForShoppingListItem(shoppingListItemGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": result})
}
