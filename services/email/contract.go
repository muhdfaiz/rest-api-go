package email

import "bitbucket.org/cliqers/shoppermate-api/systems"

type EmailServiceInterface interface {
	AddSubscriber(email, name string) *systems.ErrorData
	SendTemplate(inputs map[string]string) *systems.ErrorData
}
