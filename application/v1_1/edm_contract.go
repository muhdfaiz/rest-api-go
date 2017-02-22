package v1_1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// EdmServiceInterface is a contract that defines the method needed for EDM Service.
type EdmServiceInterface interface {
	SendEdmForInsufficientFunds(userGUID string, data SendEdmInsufficientFunds) *systems.ErrorData
}
