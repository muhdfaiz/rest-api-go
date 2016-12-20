package v1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
)

type DealCashbackTransactionHandler struct {
	DealCashbackTransactionService DealCashbackTransactionServiceInterface
}

func (dcth *DealCashbackTransactionHandler) Create(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	userGUID := c.Param("guid")

	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("create deal cashback transaction"))
		return
	}

	createDealCashbackTransaction := &CreateDealCashbackTransaction{}

	if err := Binding.Bind(createDealCashbackTransaction, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	receipt := c.Request.MultipartForm.File["receipt_image"]

	if len(receipt) < 1 {
		err := &systems.Error{}
		c.JSON(http.StatusUnprocessableEntity, err.FileRequireErrors("receipt_image"))
		return
	}

	dealCashbackGUIDs := createDealCashbackTransaction.DealCashbackGuids

	relations := "transactiontypes,transactionstatuses,dealcashbacktransactions,dealcashbacktransactions.dealcashbacks,dealcashbacktransactions.dealcashbacks.deals"

	result, err := dcth.DealCashbackTransactionService.CreateTransaction(receipt[0], userGUID, dealCashbackGUIDs, relations)

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
