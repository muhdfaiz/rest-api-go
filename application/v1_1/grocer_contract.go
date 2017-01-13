package v1_1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// GrocerTransformerInterface is a contract that defines the method needed for Grocer Transformer.
type GrocerTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *GrocerResponse
}

// GrocerServiceInterface is a contract that defines the method needed for Grocer Service.
type GrocerServiceInterface interface {
	CheckGrocerPublishOrNotByGUID(grocerGUID string) (*Grocer, *systems.ErrorData)
	CheckGrocerExistOrNotByGUID(grocerGUID string) (*Grocer, *systems.ErrorData)
	GetGrocerByID(grocerID int, relations string) *Grocer
	GetAllGrocers(pageNumber, pageLimit, relations string) ([]*Grocer, int)
	GetAllGrocersIncludingDeals(userGUID, latitude, longitude string) []*Grocer
}

// GrocerRepositoryInterface is a contract that defines the method needed for Grocer Repository.
type GrocerRepositoryInterface interface {
	GetAll(pageNumber string, pageLimit string, relations string) ([]*Grocer, int)
	GetAllGrocersThoseOnlyHaveDeal() []*Grocer
	GetByID(id int, relations string) *Grocer
	GetByGUID(grocerGUID, relations string) *Grocer
	GetByGUIDAndStatus(grocerGUID, status, relations string) *Grocer
}
