package v11

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemCategoryHandler will handle all request related to Item Category resource.
type ItemCategoryHandler struct {
	ItemCategoryService ItemCategoryServiceInterface
}

// ViewAll function used to retrieve all item categories
func (ich *ItemCategoryHandler) ViewAll(context *gin.Context) {
	itemCategoryNames, total := ich.ItemCategoryService.GetItemCategoryNames()

	itemCategories := ich.ItemCategoryService.TransformItemCategories(itemCategoryNames, total)

	context.JSON(http.StatusOK, itemCategories)
}

func (ich *ItemCategoryHandler) ViewGrocerCategoriesThoseHaveDealsIncludingDeals(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")
	grocerGUID := context.Param("grocer_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view grocer categories including deals"))
		return
	}

	queryStringValidationRules := map[string]string{
		"latitude":                "required,latitude",
		"longitude":               "required,longitude",
		"deal_limit_per_category": "required,numeric",
	}

	error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

	if error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	latitude := context.Query("latitude")
	longitude := context.Query("longitude")
	dealLimitPerSubcategory := context.DefaultQuery("deal_limit_per_category", "5")

	grocerCategoriesIncludingDeals, error := ich.ItemCategoryService.GetGrocerCategoriesThoseHaveDealsIncludingDeals(userGUID, grocerGUID, latitude, longitude, dealLimitPerSubcategory, "")

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": grocerCategoriesIncludingDeals})
}
