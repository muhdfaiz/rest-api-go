package v1_1

import (
	"net/http"
	"strconv"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SmsHandler Struct
type SmsHandler struct {
	UserRepository    UserRepositoryInterface
	SmsService        SmsServiceInterface
	SmsHistoryService SmsHistoryServiceInterface
	DeviceService     DeviceServiceInterface
	UserService       UserServiceInterface
}

// Send function used to send sms to the user during login & registration
func (sh *SmsHandler) Send(context *gin.Context) {
	smsData := &SmsSend{}

	if error := Binding.Bind(smsData, context); error != nil {
		context.JSON(http.StatusBadRequest, error)
		return
	}

	if smsData.Event == "register" {
		error := sh.UserService.CheckUserPhoneNumberDuplicate(smsData.RecipientNo)

		if error != nil {
			statusCode, _ := strconv.Atoi(error.Error.Status)
			context.JSON(statusCode, error)
			return
		}
	}

	debug := context.Query("debug")

	if debug == "1" {
		smsHistory := &SmsHistory{
			GUID:             Helper.GenerateUUID(),
			Provider:         "moceansms",
			Text:             "Your verification code is debug - Shoppermate",
			SmsID:            "shoppermate_debug",
			RecipientNo:      smsData.RecipientNo,
			VerificationCode: "9999",
			Event:            smsData.Event,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		context.JSON(http.StatusOK, gin.H{"data": smsHistory})
		return
	}

	error := sh.SmsHistoryService.CheckIfUserReachSmsLimitForToday(smsData.RecipientNo, smsData.Event)

	if error != nil {
		statusCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(statusCode, error)
		return
	}

	smsHistory := sh.SmsHistoryService.GetLatestSmsHistoryByPhoneNoAndEventName(smsData.RecipientNo, smsData.Event)

	if smsHistory.GUID != "" {
		error = sh.SmsHistoryService.CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsHistory.CreatedAt)

		if error != nil {
			statusCode, _ := strconv.Atoi(error.Error.Status)
			context.JSON(statusCode, error)
			return
		}
	}

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	// Send SMS verification code
	sentSmsData, error := sh.SmsService.SendVerificationCode(dbTransaction, smsData.RecipientNo, smsData.Event)

	if error != nil {
		dbTransaction.Rollback()
		statusCode, _ := strconv.Atoi(error.Error.Status)
		context.JSON(statusCode, error)
		return
	}

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": sentSmsData.(*SmsHistory)})
}

// Verify function used to verify sms verification code during login & registration
// Return JWT Token if sms verification code valid
func (sh *SmsHandler) Verify(context *gin.Context) {
	smsData := SmsVerification{}

	// Bind request based on content type and validate request data
	if error := Binding.Bind(&smsData, context); error != nil {
		context.JSON(http.StatusBadRequest, error)
		return
	}

	device := sh.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(smsData.DeviceUUID)

	if device.UserGUID == nil {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("Device", "device_uuid", smsData.DeviceUUID))
	}

	user := sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	debug := context.Query("debug")

	event := "login"

	if user.PhoneNo != smsData.PhoneNo {
		event = "update"
	}

	if debug != "1" {
		error := sh.SmsHistoryService.VerifyVerificationCode(smsData.PhoneNo, strings.ToLower(smsData.VerificationCode), event)

		if error != nil {
			dbTransaction.Rollback()
			errorCode, _ := strconv.Atoi(error.Error.Status)
			context.JSON(errorCode, error)
			return
		}
	}

	error := sh.DeviceService.ReactivateDevice(dbTransaction, device.GUID)

	if error != nil {
		dbTransaction.Rollback()
		context.JSON(http.StatusInternalServerError, error)
		return
	}

	debugToken := context.Query("debug_token")

	jwtToken, error := JWT.GenerateToken(user.GUID, smsData.PhoneNo, smsData.DeviceUUID, debugToken)

	if error != nil {
		dbTransaction.Rollback()
		context.JSON(http.StatusInternalServerError, error)
	}

	dbTransaction.Commit()

	user = sh.UserService.CalculateAllTimeAmountAndPendingAmount(user)

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	context.JSON(http.StatusOK, gin.H{"data": response})
}
