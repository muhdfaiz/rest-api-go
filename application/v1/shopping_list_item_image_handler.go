package v1

import (
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

//ShoppingListItemImageHandler used to handle all request related to shopping list item image
type ShoppingListItemImageHandler struct {
	UserRepository                  UserRepositoryInterface
	ShoppingListRepository          ShoppingListRepositoryInterface
	ShoppingListItemRepository      ShoppingListItemRepositoryInterface
	ShoppingListItemImageService    ShoppingListItemImageServiceInterface
	ShoppingListItemImageFactory    ShoppingListItemImageFactoryInterface
	ShoppingListItemImageRepository ShoppingListItemImageRepositoryInterface
}

// View function used to retrieve shopping list item image from database
func (sliih *ShoppingListItemImageHandler) View(c *gin.Context) {
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

	// Retrieve shopping list by guid and user guid
	shoppingList := sliih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := sliih.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemImageGUID := c.Param("image_guid")

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve shopping list item image by item guid and guid
	shoppingListItemImage := sliih.ShoppingListItemImageRepository.GetByItemGUIDAndGUID(shoppingListItemGUID, shoppingListItemImageGUID, relations)

	// If shopping list item image GUID empty return error message
	if shoppingListItemImage.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item Image", "guid", shoppingListItemImageGUID))
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": shoppingListItemImage})
}

// Create function used to store shopping list item image in database
func (sliih *ShoppingListItemImageHandler) Create(c *gin.Context) {
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
	shoppingList := sliih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := sliih.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	createImageData := &CreateImage{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(createImageData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve images in post body
	images := c.Request.MultipartForm.File["images"]

	// If shopping list item images not exist in post data return error message
	if len(images) < 1 {
		err := &systems.Error{}
		c.JSON(http.StatusUnprocessableEntity, err.FileRequireErrors("images"))
		return
	}

	err := &systems.ErrorData{}

	// Upload shopping list item images
	uploadedImages, err := sliih.ShoppingListItemImageService.UploadImages(images)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Store uploaded shopping list item image into database
	result, err := sliih.ShoppingListItemImageFactory.Create(userGUID, shoppingListGUID, shoppingListItemGUID, uploadedImages)

	//Return error message if failed to store uploaded shopping list item image into database
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})

}

// Delete function used to delete multiple shopping list item images
func (sliih *ShoppingListItemImageHandler) Delete(c *gin.Context) {
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

	// Retrieve shopping list by guid and user guid
	shoppingList := sliih.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list GUID empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Retrieve shopping list item guid in url
	shoppingListItemGUID := c.Param("item_guid")

	// Retrieve shopping list item by guid
	shoppingListItem := sliih.ShoppingListItemRepository.GetByGUID(shoppingListItemGUID, "")

	// If shopping list item GUID empty return error message
	if shoppingListItem.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item", "guid", shoppingListItemGUID))
		return
	}

	// Retrieve shopping list item image guid in url
	shoppingListItemImageGUIDs := c.Param("image_guids")

	// Split on comma.
	splitShoppingListItemImageGUID := strings.Split(shoppingListItemImageGUIDs, ",")

	shoppingListItemImageURLs := make([]string, len(splitShoppingListItemImageGUID))

	for key, shoppingListItemImageGUID := range splitShoppingListItemImageGUID {
		// Retrieve shopping list item image by item guid and guid
		shoppingListItemImage := sliih.ShoppingListItemImageRepository.GetByItemGUIDAndGUID(shoppingListItemGUID, shoppingListItemImageGUID, "")

		// If shopping list item image GUID empty return error message
		if shoppingListItemImage.GUID == "" {
			c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List Item Image", "guid", shoppingListItemImageGUID))
			return
		}

		shoppingListItemImageURLs[key] = shoppingListItemImage.URL
	}

	// Delete shopping list item image from database
	err := sliih.ShoppingListItemImageFactory.Delete("guid", splitShoppingListItemImageGUID, shoppingListItemImageURLs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted shopping list item image with GUIDs " + shoppingListItemImageGUIDs

	c.JSON(http.StatusOK, gin.H{"data": result})
}
