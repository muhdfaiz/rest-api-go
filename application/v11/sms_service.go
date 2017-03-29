package v11

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jinzhu/gorm"

	"os"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// SmsService used to handle application logic related to SMS resource.
type SmsService struct {
	DB                   *gorm.DB
	SmsHistoryRepository SmsHistoryRepositoryInterface
}

// SendVerificationCode function handle sending sms contain verification code during registration & login
func (sf *SmsService) SendVerificationCode(dbTransaction *gorm.DB, phoneNo, eventName string) (interface{}, *systems.ErrorData) {
	smsVerificationCode := Helper.RandomString("Digit", 4, "", "")

	smsText := fmt.Sprintf("Your verification code is %s - Shoppermate", smsVerificationCode)

	smsHistory := make(map[string]string)

	if os.Getenv("SEND_SMS") == "true" {
		smsResponse, err := sf.send(smsText, phoneNo)

		if err != nil {
			return nil, err
		}

		if smsResponse == nil || smsResponse["status"] == "failed" {
			return nil, Error.InternalServerError(smsResponse["message"], systems.FailedToSendSMS)
		}

		smsHistory["sms_id"] = smsResponse["sms_id"]
	}

	smsHistory["guid"] = Helper.GenerateUUID()
	smsHistory["provider"] = "moceansms"
	smsHistory["text"] = smsText
	smsHistory["recipient_no"] = phoneNo
	smsHistory["verification_code"] = smsVerificationCode
	smsHistory["event"] = eventName

	result, err := sf.SmsHistoryRepository.Create(dbTransaction, smsHistory)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Send SMS message
func (sf *SmsService) send(message, recipientNumber string) (map[string]string, *systems.ErrorData) {
	apiURL := os.Getenv("MOCEAN_SMS_URL")
	username := os.Getenv("MOCEAN_SMS_USERNAME")
	password := os.Getenv("MOCEAN_SMS_PASSWORD")
	coding := os.Getenv("MOCEAN_SMS_CODING")
	subject := os.Getenv("MCOEAN_SMS_SUBJECT")

	// Send request to SMS Server for sending sms
	req, err := http.NewRequest("GET", apiURL, nil)

	if err != nil {
		return nil, Error.InternalServerError(err.Error(), systems.FailedToSendSMS)
	}

	req.Header.Add("Connection", "close")

	// Append Query String Parameters to the URL become like this
	// http://183.81.161.84:13016/cgi-bin/sendsms?username=yourusername&password=yourpassword&from=subject&to=601234567&coding=1&text=Hello
	q := req.URL.Query()
	q.Add("username", username)
	q.Add("password", password)
	q.Add("from", subject)
	q.Add("to", recipientNumber)
	q.Add("coding", coding)
	q.Add("text", message)

	req.URL.RawQuery = q.Encode()

	resp, err2 := http.Get(req.URL.String())

	if err2 != nil {
		return nil, Error.InternalServerError(err2.Error(), systems.FailedToSendSMS)
	}

	defer resp.Body.Close()

	// Get Response Body
	response, err2 := ioutil.ReadAll(resp.Body)

	if err2 != nil {
		return nil, Error.InternalServerError(err2.Error(), systems.FailedToSendSMS)
	}

	parsedResponse, _ := url.ParseQuery(string(response))

	result := make(map[string]string)

	if parsedResponse["status"][0] == "0" {
		result["status"] = "success"
		result["sms_id"] = parsedResponse["msgid"][0]
	} else {
		result["status"] = "failed"
		result["message"] = parsedResponse["err_msg"][0]
	}

	result["status_code"] = parsedResponse["status"][0]

	return result, nil

}
