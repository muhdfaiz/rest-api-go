package v1_1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// EdmHandler will handle all task to send EDM template for various event.
// API will request email API at this URL
// http://api.shoppermate.com:5000/email/send-template
type EdmHandler struct {
	EdmService EdmServiceInterface
}

// InsufficientFunds function used to send EDM when user trying to cashout and the amount of
// cashout available is below minimum cashout amount.
func (eh *EdmHandler) InsufficientFunds(context *gin.Context) {
	userGUID := context.Param("guid")

	userToken := context.MustGet("Token").(map[string]string)

	if userToken["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("Send EDM"))
		return
	}

	edmData := SendEdmInsufficientFunds{}

	if error := Binding.Bind(&edmData, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	error := eh.EdmService.SendEdmForInsufficientFunds(dbTransaction, userGUID, edmData)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	result := make(map[string]string)
	result["message"] = "Successfully send EDM insufficient funds to user with GUID " + userGUID

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": result})
}
