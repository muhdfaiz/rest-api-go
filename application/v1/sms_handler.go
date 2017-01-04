package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SmsHandler Struct
type SmsHandler struct {
	UserRepository       UserRepositoryInterface
	SmsService           SmsServiceInterface
	SmsHistoryRepository SmsHistoryRepositoryInterface
	DeviceService        DeviceServiceInterface
}

// Send function used to send sms to the user during login & registration
func (sh *SmsHandler) Send(context *gin.Context) {
	smsData := &SmsSend{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(smsData, context); err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by GUID
	user := sh.UserRepository.GetByGUID(smsData.UserGUID, "")

	// If user GUID empty return error message
	if user.GUID == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", smsData.UserGUID))
		return
	}

	debug := context.Query("debug")

	if debug == "1" {
		smsHistory := &SmsHistory{
			GUID:             Helper.GenerateUUID(),
			UserGUID:         smsData.UserGUID,
			Provider:         "moceansms",
			Text:             "Your verification code is debug - Shoppermate",
			SmsID:            "shoppermate_debug",
			RecipientNo:      smsData.RecipientNo,
			VerificationCode: "9999",
			Status:           "0",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		context.JSON(http.StatusOK, gin.H{"data": smsHistory})
		return
	}
	// Retrieve sms history by recipient_nos
	smsHistory := sh.SmsHistoryRepository.GetByRecipientNo(smsData.RecipientNo)

	// // If recipient_no not empty
	if smsHistory.GUID != "" {
		// Calculate time interval in second between current time and last sms sent time
		interval := sh.SmsHistoryRepository.CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsHistory.CreatedAt)

		// If time interval in second below 250 return error message
		if interval < 250 {
			durationUserMustWait := 250 - interval
			errorMesg := Error.GenericError("500", systems.FailedToSendSMS, systems.TitleSentSmsError,
				"", fmt.Sprintf(systems.ErrorSentSms, strconv.Itoa(durationUserMustWait)))
			context.JSON(http.StatusBadRequest, errorMesg)
			return
		}
	}

	// Check if user already request 3 sms per day
	// Return error if user already reached 3 sms per day
	// todayDate := time.Now().UTC().Format("2006-01-02")

	// type Count struct {
	// 	Count string
	// }
	// row := DB.Table("sms_histories").Where("recipient_no = ? AND date(created_at) = ?", smsData.RecipientNo, todayDate).Select("count(*)")
	// fmt.Println(row.Scan(&Count{}))
	// os.Exit(0)

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	// Send SMS verification code
	sentSmsData, err := sh.SmsService.SendVerificationCode(dbTransaction, smsData.RecipientNo, smsData.UserGUID)

	if err != nil {
		dbTransaction.Rollback()
		statusCode, _ := strconv.Atoi(err.Error.Status)
		context.JSON(statusCode, err)
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
	if err := Binding.Bind(&smsData, context); err != nil {
		context.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by phone no
	user := sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	// If user phone_no empty return error message
	if user.PhoneNo == "" {
		context.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "phone_no", smsData.PhoneNo))
		return
	}

	// Retrieve device by uuid
	device := sh.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(smsData.DeviceUUID)

	dbTransaction := context.MustGet("DB").(*gorm.DB).Begin()

	// If Device User GUID empty, update device with User GUID
	if device.UserGUID == nil {
		_, err := sh.DeviceService.UpdateDevice(dbTransaction, smsData.DeviceUUID, UpdateDevice{UserGUID: user.GUID})

		if err != nil {
			dbTransaction.Rollback()
			context.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	debug := context.Query("debug")

	if debug != "1" {
		// Verify Sms verification code
		smsHistory := sh.SmsHistoryRepository.VerifyVerificationCode(smsData.PhoneNo, strings.ToLower(smsData.VerificationCode))

		// If sms history record not found return error message
		if smsHistory == nil {
			dbTransaction.Rollback()
			errorMesg := Error.GenericError(strconv.Itoa(http.StatusBadRequest), systems.VerificationCodeInvalid,
				systems.TitleVerificationCodeInvalid, "", fmt.Sprintf(systems.ErrorVerificationCodeInvalid, smsData.VerificationCode))
			context.JSON(http.StatusBadRequest, errorMesg)
			return
		}
	}

	// Set user status to verified
	err := sh.UserRepository.Update(dbTransaction, user.GUID, map[string]interface{}{"verified": 1})

	if err != nil {
		dbTransaction.Rollback()
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	// Set deleted_at column in devices table to null
	err = sh.DeviceService.ReactivateDevice(dbTransaction, device.GUID)

	if err != nil {
		dbTransaction.Rollback()
		context.JSON(http.StatusInternalServerError, err)
		return
	}

	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateToken(user.GUID, smsData.PhoneNo, smsData.DeviceUUID)

	if err != nil {
		dbTransaction.Rollback()
		context.JSON(http.StatusInternalServerError, err)
	}

	// Retrieve user by phone no
	user = sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	dbTransaction.Commit()

	context.JSON(http.StatusOK, gin.H{"data": response})
}
