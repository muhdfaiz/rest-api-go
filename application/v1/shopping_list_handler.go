package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ShoppingListHandler will handle all request related to user shopping list
type ShoppingListHandler struct {
	UserRepository          UserRepositoryInterface
	OccasionRepository      OccasionRepositoryInterface
	ShoppingListFactory     ShoppingListFactoryInterface
	ShoppingListRepository  ShoppingListRepositoryInterface
	ShoppingListItemFactory ShoppingListItemFactoryInterface
}

// View function used to retrieve User Shopping List
func (slh *ShoppingListHandler) View(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view shopping list"))
		return
	}

	// Retrieve query string for relations
	relations := c.DefaultQuery("include", "")

	// Retrieve User Shopping List by User GUID
	shoppingLists := slh.ShoppingListRepository.GetByUserGUID(userGUID, relations)

	if len(shoppingLists) <= 0 {
		sampleShoppingList := CreateShoppingList{OccasionGUID: "714a61a9-0aaa-5af4-86fe-0c8967463270", Name: "My Family's Grocery List"}
		slh.ShoppingListFactory.Create(userGUID, sampleShoppingList)

		// Retrieve Sample User Shopping List
		shoppingLists = slh.ShoppingListRepository.GetByUserGUID(userGUID, "Occasions")
	}

	c.JSON(http.StatusOK, gin.H{"data": shoppingLists})

}

// Create function used to create user shopping list
func (slh *ShoppingListHandler) Create(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	shoppingListData := CreateShoppingList{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&shoppingListData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve occasion by guid
	occasion := slh.OccasionRepository.GetByGUID(shoppingListData.OccasionGUID)

	// If Occasion GUID empty return error message
	if occasion.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Occasion", "guid", shoppingListData.OccasionGUID))
		return
	}

	// Retrieve user shopping list by User GUID and Shopping List name
	shoppingList := slh.ShoppingListRepository.GetByUserGUIDOccasionGUIDAndName(userGUID, shoppingListData.Name, shoppingListData.OccasionGUID, "")

	// If Shopping List name already exist return error message
	if shoppingList.Name != "" {
		c.JSON(http.StatusConflict, Error.DuplicateValueErrors("Shopping List", "name", shoppingListData.Name))
		return
	}

	createdShoppingList, err := slh.ShoppingListFactory.Create(userGUID, shoppingListData)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	result := slh.ShoppingListRepository.GetByGUID(createdShoppingList.GUID, "")

	c.JSON(http.StatusOK, gin.H{"data": result})

}

// Update function used to update user shopping list
func (slh *ShoppingListHandler) Update(c *gin.Context) {
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

	// Retrieve shopping list by Shopping List GUID
	shoppingList := slh.ShoppingListRepository.GetByGUIDAndUserGUID(shoppingListGUID, userGUID, "")

	// If shopping list guid empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	shoppingListData := UpdateShoppingList{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&shoppingListData, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve occasion by guid
	if shoppingListData.OccasionGUID != "" {
		occasion := slh.OccasionRepository.GetByGUID(shoppingListData.OccasionGUID)

		// If Occasion GUID empty return error message
		if occasion.GUID == "" {
			c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Occasion", "guid", shoppingListData.OccasionGUID))
			return
		}
	}

	if shoppingListData.Name != "" {
		// Retrieve user shopping list by User GUID and Shopping List name
		shoppingList := slh.ShoppingListRepository.GetByUserGUIDOccasionGUIDAndName(userGUID, shoppingListData.Name, shoppingListData.OccasionGUID, "")

		// If Shopping List name already exist return error message
		if shoppingList.Name != "" && shoppingList.GUID != shoppingListGUID {
			c.JSON(http.StatusConflict, Error.DuplicateValueErrors("Shopping List", "name", shoppingListData.Name))
			return
		}
	}

	err := slh.ShoppingListFactory.Update(userGUID, shoppingListGUID, shoppingListData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	result := slh.ShoppingListRepository.GetByGUID(shoppingListGUID, "")

	c.JSON(http.StatusOK, gin.H{"data": result})

}

// Delete function used to soft delete device by setting current timeo the deleted_at column
func (slh *ShoppingListHandler) Delete(c *gin.Context) {
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

	// Retrieve shopping list by Shopping List GUID
	shoppingList := slh.ShoppingListRepository.GetByGUID(shoppingListGUID, "")

	// If shopping list guid empty return error message
	if shoppingList.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Shopping List", "guid", shoppingListGUID))
		return
	}

	// Soft delete shopping list
	err := slh.ShoppingListFactory.Delete("guid", shoppingListGUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Soft delete Shopping List Item incuding relationship
	err = slh.ShoppingListItemFactory.DeleteByShoppingListGUID(shoppingList.GUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Response data
	result := make(map[string]string)
	result["message"] = "Successfully deleted shopping list with guid " + shoppingList.GUID

	c.JSON(http.StatusOK, gin.H{"data": result})
}
