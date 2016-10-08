package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SmsHandler Struct
type SmsHandler struct {
	UserRepository       UserRepositoryInterface
	UserFactory          UserFactoryInterface
	SmsService           SmsServiceInterface
	SmsHistoryRepository SmsHistoryRepositoryInterface
	DeviceRepository     DeviceRepositoryInterface
	DeviceFactory        DeviceFactoryInterface
}

// Send function used to send sms to the user during login & registration
func (sh *SmsHandler) Send(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()
	smsData := &SmsSend{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(smsData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by GUID
	user := sh.UserRepository.GetByGUID(smsData.UserGUID, "")

	// If user GUID empty return error message
	if user.GUID == "" {
		DB.Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "guid", smsData.UserGUID))
		return
	}

	// Retrieve sms history by recipient_nos
	smsHistory := sh.SmsHistoryRepository.GetByRecipientNo(smsData.RecipientNo)

	// If recipient_no not empty
	if smsHistory.RecipientNo != "" {
		// Calculate time interval in second between current time and last sms sent time
		interval := sh.SmsHistoryRepository.CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsHistory.CreatedAt)

		// If time interval in second below 250 return error message
		if interval < 250 {
			DB.Close()
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
		DB.Rollback().Close()
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	DB.Commit().Close()
	c.JSON(http.StatusOK, gin.H{"data": sentSmsData.(*SmsHistory)})
}

// Verify function used to verify sms verification code during login & registration
// Return JWT Token if sms verification code valid
func (sh *SmsHandler) Verify(c *gin.Context) {
	DB := c.MustGet("DB").(*gorm.DB).Begin()

	smsData := SmsVerification{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&smsData, c); err != nil {
		DB.Close()
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Retrieve user by phone no
	user := sh.UserRepository.GetByPhoneNo(smsData.PhoneNo, "")

	// If user phone_no empty return error message
	if user.PhoneNo == "" {
		DB.Close()
		c.JSON(http.StatusBadRequest, Error.ResourceNotFoundError("User", "phone_no", smsData.PhoneNo))
		return
	}

	// Retrieve device by uuid
	device := sh.DeviceRepository.GetByUUIDUnscoped(smsData.DeviceUUID)

	// If Device User GUID empty, update device with User GUID
	if device.UserGUID == "" {
		err := sh.DeviceFactory.Update(smsData.DeviceUUID, UpdateDevice{UserGUID: user.GUID})

		if err != nil {
			DB.Rollback().Close()
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	DB.Commit()
	DB = c.MustGet("DB").(*gorm.DB).Begin()

	// Verify Sms verification code
	smsHistory := sh.SmsHistoryRepository.VerifyVerificationCode(smsData.PhoneNo, strings.ToLower(smsData.VerificationCode))

	// If sms history record not found return error message
	if smsHistory == nil {
		DB.Close()
		errorMesg := Error.GenericError(strconv.Itoa(http.StatusBadRequest), systems.VerificationCodeInvalid,
			systems.TitleVerificationCodeInvalid, "", fmt.Sprintf(systems.ErrorVerificationCodeInvalid, smsData.VerificationCode))
		c.JSON(http.StatusBadRequest, errorMesg)
		return
	}

	// Set user status to verified
	err := sh.UserFactory.Update(smsHistory.UserGUID, map[string]interface{}{"verified": 1})

	if err != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Set deleted_at column in devices table to null
	result := DB.Unscoped().Model(&Device{}).Update("deleted_at", nil)

	if result.Error != nil {
		DB.Rollback().Close()
		c.JSON(http.StatusInternalServerError, Error.InternalServerError(result.Error, systems.DatabaseError))
		return
	}

	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateToken(smsHistory.UserGUID, smsHistory.RecipientNo, smsData.DeviceUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	DB.Commit().Close()

	response := make(map[string]interface{})
	response["user"] = user
	response["access_token"] = jwtToken

	c.JSON(http.StatusOK, gin.H{"data": response})
}
