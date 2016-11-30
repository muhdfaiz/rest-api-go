package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CashoutTransactionHandler struct {
	CashoutTransactionService CashoutTransactionServiceInterface
}

func (cth CashoutTransactionHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create cashout transaction"))
		return
	}

	createCashoutTransaction := &CreateCashoutTransaction{}

	if err := Binding.Bind(createCashoutTransaction, context); err != nil {
		context.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	transaction, error := cth.CashoutTransactionService.CreateCashoutTransaction(userGUID, createCashoutTransaction)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}
    
	context.JSON(http.StatusOK, gin.H{"data": transaction})
}
