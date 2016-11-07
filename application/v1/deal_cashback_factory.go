package v1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

type DealCashbackFactoryInterface interface {
	Create(userGUID string, data CreateDealCashback) (*DealCashback, *systems.ErrorData)
	DeleteByUserGUIDAndDealGUID(userGUID string, dealGUID string) *systems.ErrorData
}

// DealCashbackFactory will handle all task related to create, update and delete user deal cashback
type DealCashbackFactory struct {
	DB *gorm.DB
}

// Create function used to create deal cashback for user
func (dcf *DealCashbackFactory) Create(userGUID string, data CreateDealCashback) (*DealCashback, *systems.ErrorData) {
	dealCashback := &DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		DealGUID:         data.DealGUID,
	}

	result := dcf.DB.Create(dealCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*DealCashback), nil
}

func (dcf *DealCashbackFactory) DeleteByUserGUIDAndDealGUID(userGUID string, dealGUID string) *systems.ErrorData {
	result := dcf.DB.Where("user_guid = ? AND deal_guid = ?", userGUID, dealGUID).Delete(&DealCashback{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}
