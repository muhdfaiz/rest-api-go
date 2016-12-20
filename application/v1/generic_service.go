package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"
import "strconv"

type GenericService struct {
	GenericRepository GenericRepositoryInterface
}

// CheckGenericExistOrNotByID function used to check if generic category exist or not by retrieve
// generic category by ID through Generic Repository.
func (gs *GenericService) CheckGenericExistOrNotByID(genericID int) (*Generic, *systems.ErrorData) {
	generic := gs.GenericRepository.GetByID(genericID, "")

	genericIDInString := strconv.Itoa(genericID)

	if generic.ID == 0 {
		return nil, Error.ResourceNotFoundError("Generic", "id", genericIDInString)
	}

	return generic, nil
}

// GetAllGeneric function used to retrieve all generic category through Generic Repository.
func (gs *GenericService) GetAllGeneric(pageNumber, pageLimit, relations string) ([]*Generic, int) {
	generics, totalGeneric := gs.GenericRepository.GetAll(pageNumber, pageLimit, relations)

	return generics, totalGeneric
}

func (gs *GenericService) GetLatestUpdate(lastSyncDate, pageNumber, pageLimit, relations string) ([]*Generic, int) {
	generics, totalGeneric := gs.GenericRepository.GetByUpdatedAtGreaterThanLastSyncDate(lastSyncDate, pageNumber, pageLimit, relations)

	return generics, totalGeneric
}
