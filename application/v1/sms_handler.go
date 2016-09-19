package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/shoppermate/systems"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// SmsHandler Struct
type SmsHandler struct{}

// Send function used to send sms to the user during login & registration
func (sh *SmsHandler) Send(c *gin.Context) {
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	smsData := &SmsSend{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(smsData, c); err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	// Check UserGUID valid or not.
	// Return error if not valid
	userRepository := &UserRepository{DB: tx}
	user := userRepository.GetByGUID(smsData.UserGUID)
	if user.GUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "guid", smsData.UserGUID)))
		return
	}

	// Check LastSmsSent interval. If interval below 250 return error message
	smsHistoryRepository := &SmsHistoryRepository{DB: tx}
	smsHistory := smsHistoryRepository.GetByRecipientNo(smsData.RecipientNo)

	if smsHistory.RecipientNo != "" {
		// Calculate time interval in second between current time and last sms sent time
		interval := smsHistoryRepository.CalculateIntervalBetweenCurrentTimeAndLastSmsSentTime(smsHistory.CreatedAt)

		// If time interval in second below 250 return error message
		if interval < 250 {
			durationUserMustWait := 250 - interval
			errorMesg := ErrorMesg.GenericError("500", systems.FailedToSendSMS, systems.TitleSentSmsError,
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

	smsService := &SmsService{DB: tx}
	sentSmsData, err := smsService.SendVerificationCode("60174862127", smsData.UserGUID)

	if err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": sentSmsData.(*SmsHistory)})
}

// Verify function used to verify sms verification code during login & registration
// Return JWT Token if sms verification code valid
func (sh *SmsHandler) Verify(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	db := c.MustGet("DB").(*gorm.DB)
	tx := db.Begin()

	smsData := SmsVerification{}

	// Bind request based on content type and validate request data
	if err := Binding.Bind(&smsData, c); err != nil {
		statusCode, _ := strconv.Atoi(err.Error.Status)
		c.JSON(statusCode, err)
		return
	}

	// Check Phone No valid or not.
	// Return error if not valid
	user := tx.Where(&User{PhoneNo: smsData.PhoneNo}).First(&User{}).Value.(*User)
	if user.PhoneNo == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "User"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "User", "phone_no", smsData.PhoneNo)))
		return
	}

	// Check device uuid exist or not
	// If not exist display error message
	deviceRepository := &DeviceRepository{DB: tx}
	device := deviceRepository.GetByUUIDAndUserGUIDUnscoped(smsData.DeviceUUID, user.GUID)

	// Return error message if device uuid not exist
	if device.UUID == "" {
		c.JSON(http.StatusBadRequest, ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.ResourceNotFound,
			fmt.Sprintf(systems.TitleResourceNotFoundError, "Device"), "message",
			fmt.Sprintf(systems.ErrorResourceNotFound, "Device", "uuid", smsData.DeviceUUID)))
		return
	}

	// Verify Sms verification code
	smsRepository := &SmsHistoryRepository{DB: tx}
	smsHistory := smsRepository.VerifyVerificationCode(smsData.PhoneNo, strings.ToLower(smsData.VerificationCode))
	if smsHistory == nil {
		errorMesg := ErrorMesg.GenericError(strconv.Itoa(http.StatusBadRequest), systems.VerificationCodeInvalid,
			systems.TitleVerificationCodeInvalid, "", fmt.Sprintf(systems.ErrorVerificationCodeInvalid, smsData.VerificationCode))
		c.JSON(http.StatusBadRequest, errorMesg)
		return
	}

	// Set user status to verified
	userFactory := UserFactory{DB: tx}
	err := userFactory.Update(smsHistory.UserGUID, UpdateUser{Verified: "1"})

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	result := tx.Unscoped().Model(&Device{}).Update("deleted_at", nil)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, ErrorMesg.InternalServerError(result.Error, systems.DatabaseError))
		return
	}

	jwt := &systems.Jwt{}
	jwtToken, err := jwt.GenerateJWTToken(smsHistory.UserGUID, smsHistory.RecipientNo, smsData.DeviceUUID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"data": jwtToken})
}
