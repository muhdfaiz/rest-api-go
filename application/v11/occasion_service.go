package v11

import "bitbucket.org/cliqers/shoppermate-api/systems"

// OccasionService used to handle application logic related to Occasion resource.
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

// GetLatestActiveOccasionAfterLastSyncDate function used to retrieve latest active occasion after
// last sync date.
func (os *OccasionService) GetLatestActiveOccasionAfterLastSyncDate(lastSyncDate string) *OccasionResponse {
	occasions, totalOccasion := os.OccasionRepository.GetLatestUpdateWithActiveStatus(lastSyncDate)

	occasionsData := os.OccasionTransformer.TransformCollection(occasions, totalOccasion)

	return occasionsData
}

// GetAllActiveOccasions function used to retrieve all active occasions available from database.
func (os *OccasionService) GetAllActiveOccasions() *OccasionResponse {
	occasions, totalOccasion := os.OccasionRepository.GetAllWithActiveStatus()

	occasionsData := os.OccasionTransformer.TransformCollection(occasions, totalOccasion)

	return occasionsData
}
