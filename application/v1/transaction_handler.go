package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	TransactionService TransactionServiceInterface
}

func (th *TransactionHandler) ViewDealCashbackTransaction(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	transactionGUID := context.Param("transaction_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	transaction, error := th.TransactionService.ViewDealCashbackTransactionAndUpdateReadStatus(userGUID, transactionGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": transaction})
}

func (th *TransactionHandler) ViewCashoutTransaction(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	transactionGUID := context.Param("transaction_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	transaction, error := th.TransactionService.ViewCashoutTransactionAndUpdateReadStatus(userGUID, transactionGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": transaction})
}

func (th *TransactionHandler) ViewUserTransactions(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("update shopping list"))
		return
	}

	transactionStatus := context.Query("transaction_status")

	relations := context.Query("include")
	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	isRead := context.Query("is_read")

	transactions := th.TransactionService.GetUserTransactionsForSpecificStatus(context.Request, userGUID, transactionStatus, isRead, pageNumber, pageLimit, relations)

	context.JSON(http.StatusOK, gin.H{"data": transactions})
}
