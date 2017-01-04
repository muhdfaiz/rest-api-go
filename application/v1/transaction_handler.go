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
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view deal cashback transaction"))
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
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view cashout transaction"))
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

func (th *TransactionHandler) ViewReferralCashbackTransaction(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	transactionGUID := context.Param("transaction_guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view referrral cashback transaction"))
		return
	}

	transaction, error := th.TransactionService.ViewReferralCashbackTransactionAndUpdateReadStatus(userGUID, transactionGUID)

	if error != nil {
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": transaction})
}

// ViewUserTransactions function used to view all types of user transactions.
func (th *TransactionHandler) ViewUserTransactions(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("view user transaction"))
		return
	}

	transactionStatus := context.Query("transaction_status")

	pageNumber := context.Query("page_number")
	pageLimit := context.Query("page_limit")
	isRead := context.Query("is_read")

	transactions := th.TransactionService.GetUserTransactions(context.Request, userGUID, transactionStatus, isRead, pageNumber, pageLimit)

	context.JSON(http.StatusOK, gin.H{"data": transactions})
}
