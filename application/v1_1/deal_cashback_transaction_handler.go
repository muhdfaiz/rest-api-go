package v1_1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type DealCashbackTransactionHandler struct {
	DealCashbackTransactionService DealCashbackTransactionServiceInterface
	TransactionService             TransactionServiceInterface
}

func (dcth *DealCashbackTransactionHandler) Create(context *gin.Context) {
	tokenData := context.MustGet("Token").(map[string]string)

	userGUID := context.Param("guid")

	if tokenData["user_guid"] != userGUID {
		context.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create deal cashback transaction"))
		return
	}

	createDealCashbackTransaction := &CreateDealCashbackTransaction{}

	if error := Binding.Bind(createDealCashbackTransaction, context); error != nil {
		context.JSON(http.StatusUnprocessableEntity, error)
		return
	}

	receipt := context.Request.MultipartForm.File["receipt_image"]

	if len(receipt) < 1 {
		error := &systems.Error{}
		context.JSON(http.StatusUnprocessableEntity, error.FileRequireErrors("receipt_image"))
		return
	}

	dealCashbackGUIDs := createDealCashbackTransaction.DealCashbackGuids

	relations := "transactiontypes,transactionstatuses,dealcashbacktransactions,dealcashbacktransactions.dealcashbacks,dealcashbacktransactions.dealcashbacks.deals"

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	transaction, error := dcth.DealCashbackTransactionService.CreateTransaction(dbTransaction, receipt[0], userGUID, dealCashbackGUIDs, relations)

	if error != nil {
		dbTransaction.Rollback()
		errorCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(errorCode, error)
		return
	}

	dbTransaction.Commit()

	transaction = dcth.TransactionService.ViewTransactionDetails(transaction.GUID, relations)

	context.JSON(http.StatusOK, gin.H{"data": transaction})
}
