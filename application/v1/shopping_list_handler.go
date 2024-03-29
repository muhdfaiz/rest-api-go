package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ShoppingListHandler will handle all request related to user shopping list
type ShoppingListHandler struct {
	ShoppingListService          ShoppingListServiceInterface
	ShoppingListItemImageService ShoppingListItemImageServiceInterface
}

// View function used to retrieve all user shopping lists.
// If user shopping list empty, API will create sample shopping list and return to the user.
func (slh *ShoppingListHandler) View(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view shopping list"))
		return
	}

	relations := context.DefaultQuery("include", "")

	shoppingLists, error := slh.ShoppingListService.GetUserShoppingLists(userGUID, relations)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": shoppingLists})
}

// Create function used to create user shopping list
func (slh *ShoppingListHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create shopping list"))
		return
	}

	createData := CreateShoppingList{}

	if error := Binding.Bind(&createData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	createdShoppingList, error := slh.ShoppingListService.CreateUserShoppingList(dbTransaction, userGUID, createData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": createdShoppingList})

}

// Update function used to update user shopping list
func (slh *ShoppingListHandler) Update(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	updateData := UpdateShoppingList{}

	if error := Binding.Bind(&updateData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	updatedShoppingList, error := slh.ShoppingListService.UpdateUserShoppingList(dbTransaction, userGUID, shoppingListGUID, updateData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": updatedShoppingList})

}

// Delete function used to soft delete device by setting current timeo the deleted_at column
func (slh *ShoppingListHandler) Delete(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("delete shopping list"))
		return
	}

	shoppingListGUID := context.Param("shopping_list_guid")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	error := slh.ShoppingListService.DeleteUserShoppingListIncludingItemsAndImages(dbTransaction, userGUID, shoppingListGUID)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	result := make(map[string]string)
	result["message"] = "Successfully deleted shopping list including with guid " + shoppingListGUID

	context.JSON(http.StatusOK, gin.H{"data": result})
}
