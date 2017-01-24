package email

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type EmailService struct{}

// AddSubscriber function used to add user into mailchimp list.
// Can view list of subscribers here at link below.
// https://us14.admin.mailchimp.com/lists/members/?id=68331
func (e EmailService) AddSubscriber(email, name string) *systems.ErrorData {
	values := url.Values{}
	values.Set("email", email)
	values.Add("name", name)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("SHOPPERMATE_EMAIL_API_URL")+"email/add-subscriber",
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	fmt.Println("body")
	fmt.Println(string(body))

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

	fmt.Println("Data")
	fmt.Println(values)

	req, err := http.NewRequest(
		"POST",
		os.Getenv("SHOPPERMATE_EMAIL_API_URL")+"email/send-template",
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Error.InternalServerError(err, systems.ErrorRequestShoppermateEmailAPI)
	}

	fmt.Println("body")
	fmt.Println(string(body))

	return nil
}
