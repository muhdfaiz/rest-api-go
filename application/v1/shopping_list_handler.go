package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ShoppingListHandler will handle all request related to user shopping list
type ShoppingListHandler struct {
	UserRepository UserRepositoryInterface
}

func (slh *ShoppingListHandler) Create(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB).Begin()

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	user := slh.UserRepository.GetByGUID(db, userGUID)

	// If user GUID empty return error message
	if user.GUID == "" {
		db.Rollback().Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	shoppingListData := CreateShoppingList{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&shoppingListData, c); err != nil {
		db.Rollback().Close()
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	db.Commit().Close()

}
