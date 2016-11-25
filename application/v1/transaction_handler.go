package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	TransactionService TransactionServiceInterface
	UserService        UserServiceInterface
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

	user := th.UserService.CheckUserExistOrNot(userGUID, "")

	if user.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", userGUID))
		return
	}

	transactions := th.TransactionService.GetUserTransactionsForSpecificStatus(context.Request, userGUID, transactionStatus, pageNumber, pageLimit, relations)

	context.JSON(http.StatusOK, gin.H{"data": transactions})
}
