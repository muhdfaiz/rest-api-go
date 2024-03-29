package v11

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ShoppingListItemHandler will handle all task related to Shopping List Item resource.
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

	relations := context.Query("include")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view shopping list item"))
		return
	}

	_, error := slih.ShoppingListService.CheckUserShoppingListExistOrNot(userGUID, shoppingListGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	shoppingListItem, error := slih.ShoppingListItemService.ViewUserShoppingListItem(userGUID, shoppingListGUID, shoppingListItemGUID, relations)

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
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view all shopping list items"))
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

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	userShoppingListItems, error := slih.ShoppingListItemService.ViewAllUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, addedToCart, latitude, longitude, relations)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": userShoppingListItems})
}

// Create function used to create user shopping list item
func (slih *ShoppingListItemHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create shopping list item"))
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

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	createdShoppingListItem, error := slih.ShoppingListItemService.CreateUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, shoppingListItemToCreate)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": createdShoppingListItem})
}

// Update function used to update one of the user shopping list item
func (slih *ShoppingListItemHandler) Update(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list item"))
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

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	updatedShoppingListItem, error := slih.ShoppingListItemService.UpdateUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, shoppingListItemGUID, shoppingListItemToUpdate, "")

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	updatedShoppingListItem = slih.ShoppingListItemService.GetShoppingListItemByGUID(shoppingListItemGUID, "")

	context.JSON(http.StatusOK, gin.H{"data": updatedShoppingListItem})
}

// UpdateAll function used to update shopping list item
func (slih *ShoppingListItemHandler) UpdateAll(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list items"))
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

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	updatedShoppingListItems, error := slih.ShoppingListItemService.UpdateAllUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, updateShoppingListItemData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	updatedShoppingListItems = slih.ShoppingListItemService.GetShoppingListItemsByUserGUIDAndShoppingListGUID(userGUID, shoppingListGUID, "")

	context.JSON(http.StatusOK, gin.H{"data": updatedShoppingListItems})
}

// DeleteAll function used to delete all shopping list items including relation like image
// or to to delete all shopping list items from cart
func (slih *ShoppingListItemHandler) DeleteAll(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("delete all shopping list items"))
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

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	result, error := slih.ShoppingListItemService.DeleteUserShoppingListItem(dbTransaction, userGUID, shoppingListGUID, deleteItemInCart)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	error = slih.ShoppingListItemImageService.DeleteImagesForShoppingList(dbTransaction, shoppingListGUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": result})
}

// Delete function used to delete user shopping list item by shopping list item GUID, user GUID and shopping list GUID.
func (slih *ShoppingListItemHandler) Delete(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("delete shopping list item"))
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")
	shoppingListItemGUID := context.Param("item_guid")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	result, error := slih.ShoppingListItemService.DeleteShoppingListItemInShoppingList(dbTransaction, shoppingListItemGUID, userGUID, shoppingListGUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	error = slih.ShoppingListItemImageService.DeleteImagesForShoppingListItem(dbTransaction, shoppingListItemGUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": result})
}
