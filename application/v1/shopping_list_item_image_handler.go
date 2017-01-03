package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//ShoppingListItemImageHandler used to handle all request related to shopping list item image resource.
type ShoppingListItemImageHandler struct {
	ShoppingListItemImageService ShoppingListItemImageServiceInterface
}

// View function used to retrieve user shopping list item image from database.
func (sliih *ShoppingListItemImageHandler) View(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	shoppingListItemGUID := context.Param("item_guid")

	shoppingListItemImageGUID := context.Param("image_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	relations := context.Query("include")

	shoppingListItemImage, error := sliih.ShoppingListItemImageService.ViewUserShoppingListItemImage(userGUID, shoppingListGUID,
		shoppingListItemGUID, shoppingListItemImageGUID, relations)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": shoppingListItemImage})
}

// Create function used to store shopping list item image in database
func (sliih *ShoppingListItemImageHandler) Create(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	shoppingListItemGUID := context.Param("item_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create shopping list item image"))
		return
	}

	createShoppingListItemImageData := &CreateShoppingListItemImage{}

	if error := Binding.Bind(createShoppingListItemImageData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	imagesToUpload := context.Request.MultipartForm.File["images"]

	createdImages, error := sliih.ShoppingListItemImageService.CreateUserShoppingListItemImage(dbTransaction, userGUID, shoppingListGUID, shoppingListItemGUID, imagesToUpload)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": createdImages})
}

// Delete function used to delete multiple shopping list item images
func (sliih *ShoppingListItemImageHandler) Delete(context *gin.Context) {
	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	shoppingListGUID := context.Param("shopping_list_guid")

	shoppingListItemGUID := context.Param("item_guid")

	shoppingListItemImageGUIDs := context.Param("image_guids")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("delete shopping list item image"))
		return
	}

	error := sliih.ShoppingListItemImageService.DeleteShoppingListItemImages(dbTransaction, userGUID, shoppingListGUID, shoppingListItemGUID, shoppingListItemImageGUIDs)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	deletedResult := make(map[string]string)
	deletedResult["message"] = "Successfully deleted shopping list item image with GUIDs " + shoppingListItemImageGUIDs

	context.JSON(http.StatusOK, gin.H{"data": deletedResult})
}
