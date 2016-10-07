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

type Binding struct {
}

func (b *Binding) Bind(obj interface{}, c *gin.Context) *ErrorData {
	binding := bind.Default(c.Request.Method, c.ContentType())

	if err := binding.Bind(c.Request, obj); err != nil {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if reflect.TypeOf(err) == reflect.TypeOf(validator.ValidationErrors{}) {
			return errorMesg.ValidationErrors(err.(validator.ValidationErrors))
		}

		if err != nil {
			return errorMesg.BindingError(err)
		}
	}

	return nil

}
