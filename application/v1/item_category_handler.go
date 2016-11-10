package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemCategoryHandler struct {
	ItemCategoryService ItemCategoryServiceInterface
}

// ViewAll function used to retrieve all item categories
func (ich *ItemCategoryHandler) ViewAll(c *gin.Context) {
	itemCategoryNames, total := ich.ItemCategoryService.GetItemCategoryNames()

	itemCategories := ich.ItemCategoryService.TransformItemCategories(itemCategoryNames, total)

	c.JSON(http.StatusOK, itemCategories)
}
