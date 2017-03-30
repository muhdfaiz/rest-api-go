package v11

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/services/email"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// CashoutTransactionHandler will handle all request related to Cashout Transaction resource.
type CashoutTransactionHandler struct {
	CashoutTransactionService CashoutTransactionServiceInterface
	TransactionService        TransactionServiceInterface
	EmailService              email.EmailServiceInterface
}

// Create new cashout transaction and store in database.
// First, it will retrieve access token in context. Incoming requests to a server should create a Context.
// See auth middleware (middlewre/auth.go) how API store the access token in context when the request coming.
// Then it will grab user GUID in request URI and check if user GUIDin request URI same with user GUID in access token.
// If user GUID not same, It will return an error.
// Then it will check if user still has pending cashout transaction. It will return an error if user still has pending cashout transaction.
// Then it will bind request body to struct and validate request data.
// Lastly, it will create cashout transaction through cashout transaction service and return the newly created cashout transaction.
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
