package v1_1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/services/email"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type CashoutTransactionHandler struct {
	CashoutTransactionService CashoutTransactionServiceInterface
	TransactionService        TransactionServiceInterface
	EmailService              email.EmailServiceInterface
}

func (cth CashoutTransactionHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create cashout transaction"))
		return
	}

	error := cth.TransactionService.CheckIfUserHasPendingCashoutTransaction(userGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	createCashoutTransaction := &CreateCashoutTransaction{}

	if err := Binding.Bind(createCashoutTransaction, context); err != nil {
		context.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	transaction, error := cth.CashoutTransactionService.CreateCashoutTransaction(dbTransaction, userGUID, createCashoutTransaction)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	relations := "transactiontypes,transactionstatuses,cashouttransactions,users"

	transaction = cth.TransactionService.ViewTransactionDetails(transaction.GUID, relations)

	context.JSON(http.StatusOK, gin.H{"data": transaction})
}
