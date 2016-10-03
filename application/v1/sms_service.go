package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jinzhu/gorm"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type SmsServiceInterface interface {
	SendVerificationCode(phoneNo string, userGUID string) (interface{}, *systems.ErrorData)
	Send(message string, recipientNumber string) (map[string]string, *systems.ErrorData)
}

type SmsService struct {
	DB *gorm.DB
}

// SendVerificationCode function handle sending sms contain verification code during registration & login
func (sf *SmsService) SendVerificationCode(phoneNo string, userGUID string) (interface{}, *systems.ErrorData) {
	// Generate randomverification code 6 character (lower & digit)
	smsVerificationCode := Helper.RandomString("Digit", 4, "", "")

	// Set Sms Text
	smsText := fmt.Sprintf("Your verification code is %s - Shoppermate", smsVerificationCode)

	// Send Sms
	smsResponse, err := sf.Send(smsText, phoneNo)

	if err != nil {
		return nil, err
	}

	if smsResponse == nil || smsResponse["status"] == "failed" {
		return nil, Error.InternalServerError(smsResponse["message"], systems.FailedToSendSMS)
	}

	// Store SMS History
	m := make(map[string]string)
	m["guid"] = Helper.GenerateUUID()
	m["user_guid"] = userGUID
	m["provider"] = "moceansms"
	m["sms_id"] = smsResponse["sms_id"]
	m["text"] = smsText
	m["recipient_no"] = phoneNo
	m["verification_code"] = smsVerificationCode
	m["status"] = "0"

	smsHistoryFactory := SmsHistoryFactory{DB: sf.DB}
	sentSmsData, err := smsHistoryFactory.CreateSmsHistory(m)

	if err != nil {
		return nil, err
	}

	return sentSmsData, nil
}

// Send SMS message
func (sf *SmsService) Send(message string, recipientNumber string) (map[string]string, *systems.ErrorData) {
	apiURL := Config.Get("sms.yaml", "mocean_sms_url", "http://183.81.161.84:13016/cgi-bin/sendsms")
	username := Config.Get("sms.yaml", "mocean_sms_username", "shoppermate-api")
	password := Config.Get("sms.yaml", "mocean_sms_password", "s28Dua3p")
	coding := Config.Get("sms.yaml", "mocean_sms_coding", "1")
	subject := Config.Get("sms.yaml", "mocean_sms_subject", "Shopper Mate")

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
