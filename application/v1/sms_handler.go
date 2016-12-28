package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"strings"

	"github.com/gin-gonic/gin"
)

// SmsHandler Struct
type SmsHandler struct {
	UserRepository       UserRepositoryInterface
	SmsService           SmsServiceInterface
	SmsHistoryRepository SmsHistoryRepositoryInterface
	DeviceService        DeviceServiceInterface
}

// Send function used to send sms to the user during login & registration
func (sh *SmsHandler) Send(c *gin.Context) {
	smsData := &SmsSend{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(smsData, c); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by GUID
	user := sh.UserRepository.GetByGUID(smsData.UserGUID, "")

	// If user GUID empty return error message
	if user.GUID == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "guid", smsData.UserGUID))
		return
	}

	debug := c.Query("debug")

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

		c.JSON(http.StatusOK, gin.H{"data": smsHistory})
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
			c.JSON(http.StatusBadRequest, errorMesg)
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

	// Send SMS verification code
	sentSmsData, err := sh.SmsService.SendVerificationCode(smsData.RecipientNo, smsData.UserGUID)

	if err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sentSmsData.(*SmsHistory)})
}

// Verify function used to verify sms verification code during login & registration
// Return JWT Token if sms verification code valid
func (sh *SmsHandler) Verify(c *gin.Context) {
	smsData := SmsVerification{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&smsData, c); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by phone no
	user := sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	// If user phone_no empty return error message
	if user.PhoneNo == "" {
		c.JSON(http.StatusNotFound, Error.ResourceNotFoundError("User", "phone_no", smsData.PhoneNo))
		return
	}

	// Retrieve device by uuid
	device := sh.DeviceService.ViewDeviceByUUIDIncludingSoftDelete(smsData.DeviceUUID)

	// If Device User GUID empty, update device with User GUID
	if device.UserGUID == "" {
		_, err := sh.DeviceService.UpdateDevice(smsData.DeviceUUID, UpdateDevice{UserGUID: user.GUID})

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	debug := c.Query("debug")

	if debug != "1" {
		// Verify Sms verification code
		smsHistory := sh.SmsHistoryRepository.VerifyVerificationCode(smsData.PhoneNo, strings.ToLower(smsData.VerificationCode))

		// If sms history record not found return error message
		if smsHistory == nil {
			errorMesg := Error.GenericError(strconv.Itoa(http.StatusBadRequest), systems.VerificationCodeInvalid,
				systems.TitleVerificationCodeInvalid, "", fmt.Sprintf(systems.ErrorVerificationCodeInvalid, smsData.VerificationCode))
			c.JSON(http.StatusBadRequest, errorMesg)
			return
		}
	}

	// Set user status to verified
	err := sh.UserRepository.Update(user.GUID, map[string]interface{}{"verified": 1})

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Set deleted_at column in devices table to null
	err = sh.DeviceService.ReactivateDevice(device.GUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateToken(user.GUID, smsData.PhoneNo, smsData.DeviceUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	// Retrieve user by phone no
	user = sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	c.JSON(http.StatusOK, gin.H{"data": response})
}
