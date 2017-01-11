package v1_1

import "bitbucket.org/cliqers/shoppermate-api/systems"

// OccasionTransformerInterface is a contract that defines the method needed for Occasion Transformer.
type OccasionTransformerInterface interface {
	TransformCollection(data interface{}, totalData int) *OccasionResponse
}

// OccasionServiceInterface is a contract that defines the methods needed for Occasion Service.
type OccasionServiceInterface interface {
	CheckOccassionExistOrNot(occasionGUID string) (*Occasion, *systems.ErrorData)
	GetLatestOccasionAfterLastSyncDate(lastSyncDate string) *OccasionResponse
	GetAllOccasions() *OccasionResponse
}

// OccasionRepositoryInterface is a contract that defines the method needed for Occasion Repository.
type OccasionRepositoryInterface interface {
	GetAll() ([]*Occasion, int)
	GetLatestUpdate(lastSyncDate string) ([]*Occasion, int)
	GetByGUID(guid string) *Occasion
}
