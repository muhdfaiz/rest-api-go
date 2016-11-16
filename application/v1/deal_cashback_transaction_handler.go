package v1

import (
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
)

type DealCashbackTransactionHandler struct {
	DealCashbackTransactionService DealCashbackTransactionServiceInterface
	DealCashbackTransactionFactory DealCashbackTransactionFactoryInterface
	DealCashbackFactory            DealCashbackFactoryInterface
}

func (dcth *DealCashbackTransactionHandler) Create(c *gin.Context) {
	tokenData := c.MustGet("Token").(map[string]string)

	// Retrieve user guid in url
	userGUID := c.Param("guid")

	// If user GUID not match user GUID inside the token return error message
	if tokenData["user_guid"] != userGUID {
		c.JSON(http.StatusUnauthorized, Error.TokenIdentityNotMatchError("deal cashback transaction"))
		return
	}

	createDealCashbackTransaction := &CreateDealCashbackTransaction{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(createDealCashbackTransaction, c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve images in post body
	receipt := c.Request.MultipartForm.File["receipt_image"]

	// If shopping list item images not exist in post data return error message
	if len(receipt) < 1 {
		err := &systems.Error{}
		c.JSON(http.StatusUnprocessableEntity, err.FileRequireErrors("receipt_image"))
		return
	}

	// Upload shopping list item images
	uploadedReceipt, err := dcth.DealCashbackTransactionService.UploadReceipt(receipt[0])

	if err != nil {
		errorCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(errorCode, err)
		return
	}

	// Store uploaded shopping list item image into database
	result, err := dcth.DealCashbackTransactionFactory.Create(userGUID, uploadedReceipt)

	//Return error message if failed to store uploaded shopping list item image into database
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Retrieve shopping list item image guid in url
	dealCashbackGUIDs := c.Param("deal_cashback_guids")

	// Split on comma.
	splitDealCashbackGUID := strings.Split(dealCashbackGUIDs, ",")

	err = dcth.DealCashbackFactory.SetDealCashbackTransactionGUID(splitDealCashbackGUID, result.GUID)

	//Return error message if failed to store uploaded shopping list item image into database
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}
