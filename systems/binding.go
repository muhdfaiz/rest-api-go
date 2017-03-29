package systems

import (
	"reflect"

	"github.com/gin-gonic/gin"
	bind "github.com/gin-gonic/gin/binding"
	validator "gopkg.in/go-playground/validator.v8"
)

var (
	errorMesg = Error{}
	binding   = Binding{}
)

// Binding is a wrapper for GIN binding. API still used GIN binding to bind
// the request data into struct.
type Binding struct{}

// Bind function used to bind request data based on request header.
// It will return error if got any failed validation.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
func (b *Binding) Bind(obj interface{}, c *gin.Context) *ErrorData {
	binding := bind.Default(c.Request.Method, c.ContentType())

	if err := binding.Bind(c.Request, obj); err != nil {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Check if the error if type of go playground validator error.
		// Return error messageif matched.
		if reflect.TypeOf(err) == reflect.TypeOf(validator.ValidationErrors{}) {
			return errorMesg.ValidationErrors(err.(validator.ValidationErrors))
		}

		if err != nil {
			return errorMesg.BindingError(err)
		}
	}

	return nil

}
