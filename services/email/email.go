package email

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type EmailService struct{}

type EmailResponse struct {
	SuccessCode      int         `json:"success_code"`
	ValidationErrors interface{} `json:"validation_errors"`
	Error            interface{} `json:"error"`
}

// AddSubscriber function used to add user into mailchimp list.
// Can view list of subscribers here at link below.
// https://us14.admin.mailchimp.com/lists/members/?id=68331
func (e EmailService) AddSubscriber(email, name string) *systems.ErrorData {
	values := url.Values{}
	values.Set("email", email)
	values.Add("name", name)

	req, error := http.NewRequest(
		"POST",
		os.Getenv("SHOPPERMATE_EMAIL_API_URL")+"email/add-subscriber",
		strings.NewReader(values.Encode()),
	)

	if error != nil {
		return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, error := client.Do(req)

	if error != nil {
		return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	}

	if resp.StatusCode != 200 {
		return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	}

	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)

	fmt.Println(error)
	fmt.Println(body)
	// if error != nil {
	// 	return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	// }

	// emailResponse := EmailResponse{}

	// error = json.Unmarshal(body, &emailResponse)

	// if error != nil {
	// 	return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	// }

	// if emailResponse.SuccessCode != "200" || emailResponse.ValidationErrors != nil || emailResponse.Error != nil {
	// 	return Error.InternalServerError(error, systems.ErrorAddSubscriberToMailchimp)
	// }

	return nil
}

// SendTemplate function used to send email for different event.
// For example new user registration, cashout transaction.
func (e EmailService) SendTemplate(inputs map[string]string) *systems.ErrorData {
	values := url.Values{}

	counter := 0

	for key, input := range inputs {
		if counter == 0 {
			values.Set(key, input)
		} else {
			values.Add(key, input)
		}

		counter++
	}

	if os.Getenv("SEND_EMAIL_EVENT") == "true" {
		values.Add("env", "prod")
	} else {
		values.Add("env", "sandbox")
	}

	req, error := http.NewRequest(
		"POST",
		os.Getenv("SHOPPERMATE_EMAIL_API_URL")+"email/send-template",
		strings.NewReader(values.Encode()),
	)

	if error != nil {
		return Error.InternalServerError(error.Error(), systems.ErrorSendingEDMThroughEmailAPI)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, error := client.Do(req)

	if error != nil {
		return Error.InternalServerError(error.Error(), systems.ErrorSendingEDMThroughEmailAPI)
	}

	if resp.StatusCode != 200 {
		return Error.InternalServerError(error.Error(), systems.ErrorSendingEDMThroughEmailAPI)
	}

	defer resp.Body.Close()

	body, error := ioutil.ReadAll(resp.Body)

	if error != nil {
		return Error.InternalServerError(error.Error(), systems.ErrorSendingEDMThroughEmailAPI)
	}

	emailResponse := new(EmailResponse)

	error = json.Unmarshal(body, &emailResponse)

	if error != nil {
		return Error.InternalServerError(error.Error(), systems.ErrorSendingEDMThroughEmailAPI)
	}

	if emailResponse.SuccessCode != 200 || emailResponse.ValidationErrors != false || emailResponse.Error != nil {
		return Error.InternalServerError(error, systems.ErrorSendingEDMThroughEmailAPI)
	}

	return nil
}
