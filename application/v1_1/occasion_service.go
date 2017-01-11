package v1_1

import "bitbucket.org/cliqers/shoppermate-api/systems"

type OccasionService struct {
	OccasionRepository  OccasionRepositoryInterface
	OccasionTransformer OccasionTransformerInterface
}

// CheckOccassionExistOrNot function used to check if occasion exist or not by using occasion GUID to query database.
func (os *OccasionService) CheckOccassionExistOrNot(occasionGUID string) (*Occasion, *systems.ErrorData) {
	occasion := os.OccasionRepository.GetByGUID(occasionGUID)

	if occasion.GUID == "" {
		return nil, Error.ResourceNotFoundError("Occasion", "guid", occasionGUID)
	}

	return occasion, nil
}

// GetLatestOccasionAfterLastSyncDate function used to retrieve latest occasion after
// last sync date.
func (os *OccasionService) GetLatestOccasionAfterLastSyncDate(lastSyncDate string) *OccasionResponse {
	occasions, totalOccasion := os.OccasionRepository.GetLatestUpdate(lastSyncDate)

	occasionsData := os.OccasionTransformer.TransformCollection(occasions, totalOccasion)

	return occasionsData
}

// GetAllOccasions function used to retrieve all occasions available from database.
func (os *OccasionService) GetAllOccasions() *OccasionResponse {
	occasions, totalOccasion := os.OccasionRepository.GetAll()

	occasionsData := os.OccasionTransformer.TransformCollection(occasions, totalOccasion)

	return occasionsData
}
