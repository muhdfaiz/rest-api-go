package systems

import (
	"reflect"

	validator "gopkg.in/go-playground/validator.v8"

	"github.com/gin-gonic/gin"
)

var (
	errorMesg = Error{}
	binding   = Binding{}
)

type Binding struct {
}

func (b *Binding) Bind(obj interface{}, c *gin.Context) *ErrorData {
	err := c.Bind(obj)

	if reflect.TypeOf(err) == reflect.TypeOf(validator.ValidationErrors{}) {
		return errorMesg.ValidationErrors(err.(validator.ValidationErrors))
	}

	if err != nil {
		return errorMesg.BindingError(err)
	}

	return nil

}
