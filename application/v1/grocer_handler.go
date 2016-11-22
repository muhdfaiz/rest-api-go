package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GrocerHandler used to handle all request related to grocers
type GrocerHandler struct {
	GrocerRepository  GrocerRepositoryInterface
	GrocerTransformer GrocerTransformerInterface
}

// Index function used to retrieve all grocers
func (gh *GrocerHandler) Index(c *gin.Context) {
	// Validate query string
	err := Validation.Validate(c.Request.URL.Query(), map[string]string{"latitude": "required,latitude", "longitude": "required,longitude"})

	// If validation error return error message
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve query string in request
	offset := c.DefaultQuery("page_number", "1")
	limit := c.DefaultQuery("page_limit", "-1")
	relations := c.Query("include")
	// latitude := c.Query("latitude")
	// longitude := c.Query("longitude")

	grocers, totalGrocers := gh.GrocerRepository.GetAll(offset, limit, relations)

	result := gh.GrocerTransformer.transformCollection(c.Request, grocers, totalGrocers, limit)

	c.JSON(http.StatusOK, result)
}
