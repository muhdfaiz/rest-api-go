package v1

import (
	"net/http"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type GenericTransformerInterface interface {
	transformCollection(request *http.Request, data interface{}, totalData int, limit string) *GenericResponse
}

// GenericServiceInterface is a contract that defines the method needed for
// Generic Service.
type GenericServiceInterface interface {
	CheckGenericExistOrNotByID(genericID int) (*Generic, *systems.ErrorData)
	GetAllGeneric(pageNumber, pageLimit, relations string) ([]*Generic, int)
	GetLatestUpdate(lastSyncDate, pageNumber, pageLimit, relations string) ([]*Generic, int)
}

// GenericRepositoryInterface is acontract that defines the method needed for
// Generic Repository.
type GenericRepositoryInterface interface {
	GetAll(pageNumber, pageLimit, relations string) ([]*Generic, int)
	GetByUpdatedAtGreaterThanLastSyncDate(lastSyncDate, pageNumber, pageLimit, relations string) ([]*Generic, int)
	GetByID(genericID int, relations string) *Generic
}
