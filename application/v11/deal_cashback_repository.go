package v11

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// DealCashbackRepository will handle all CRUD function for Deal Cashback resource.
type DealCashbackRepository struct {
	BaseRepository
	DB *gorm.DB
}

// Create function used to create new deal cashback for user using user GUID and store in database.
func (dcr *DealCashbackRepository) Create(dbTransaction *gorm.DB, userGUID string, data CreateDealCashback) (*DealCashback, *systems.ErrorData) {
	dealCashback := &DealCashback{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: data.ShoppingListGUID,
		DealGUID:         data.DealGUID,
	}

	result := dbTransaction.Create(dealCashback)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*DealCashback), nil
}

// UpdateDealCashbackTransactionGUID function used to set Deal Cashback Transaction GUID for multiple deal cashback using deal cashback GUID.
// It's not update all fields available in user table but only update fields that exist in data parameter.
// Use database transaction to create new user. Don't forget to commit the transaction after used this function.
// Return nil if sucessfully updated deal cashback transaction GUID and return error if failed to update deal cashback transaction GUID.
func (dcr *DealCashbackRepository) UpdateDealCashbackTransactionGUID(dbTransaction *gorm.DB, dealCashbackGUIDs []string, dealCashbackTransactionGUID string) *systems.ErrorData {
	result := dbTransaction.Model(&DealCashback{}).Where("guid IN (?)", dealCashbackGUIDs).
		Updates(map[string]interface{}{"deal_cashback_transaction_guid": dealCashbackTransactionGUID})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByUserGUIDAndDealGUID function used to soft delete Deal Cashback using user GUID and deal GUID.
// Soft delete means it will set the current date and time to `deleted_at` field in deal cashback table.
// Use database transaction to update existing user info. Don't forget to commit the transaction after used this function.
// Return nil if sucessfully delete deal cashback GUID and return error if failed to delete deal cashback GUID.
func (dcr *DealCashbackRepository) DeleteByUserGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, dealGUID string) *systems.ErrorData {
	result := dbTransaction.Model(&DealCashback{}).Where(&DealCashback{UserGUID: userGUID, DealGUID: dealGUID}).Delete(&DealCashback{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// DeleteByUserGUIDAndShoppingListGUIDAndDealGUID function used to soft delete deal cashback using user GUID, shopping list GUID and deal GUID.
// Soft delete means it will set the current date and time to `deleted_at` field in deal cashback table.
// Use database transaction to update existing user info. Don't forget to commit the transaction after used this function.
// Return nil if sucessfully delete deal cashback GUID and return error if failed to delete deal cashback GUID.
func (dcr *DealCashbackRepository) DeleteByUserGUIDAndShoppingListGUIDAndDealGUID(dbTransaction *gorm.DB, userGUID string, shoppingListGUID string, dealGUID string) *systems.ErrorData {
	result := dbTransaction.Model(&DealCashback{}).Where(&DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, DealGUID: dealGUID}).Delete(&DealCashback{})

	if result.Error != nil {
		return Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return nil
}

// GetByGUID function used to retrieve deal cashback by deal cashback GUID.
// It will return empty deal cashback if cannot be found and return deal cashback data if matched.
func (dcr *DealCashbackRepository) GetByGUID(GUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{GUID: GUID}).Find(&dealCashback)

	return dealCashback
}

// GetByUserGUIDGroupByShoppingList function used to retrieve multiple deal cashbacks by user GUID.
func (dcr *DealCashbackRepository) GetByUserGUIDGroupByShoppingList(userGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int) {
	dealCashbacks := []*DealCashback{}

	DB := dcr.DB.Model(&DealCashback{})

	offset := dcr.SetOffsetValue(pageNumber, pageLimit)

	if relations != "" {
		DB = dcr.LoadRelations(DB, relations)
	}

	if pageLimit != "" {
		totalDealCashbacks := []*DealCashback{}

		DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
			Where(&DealCashback{UserGUID: userGUID}).Group("shopping_lists.guid").
			Find(&totalDealCashbacks)

		DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
			Where(&DealCashback{UserGUID: userGUID}).Group("shopping_lists.guid").
			Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

		return dealCashbacks, len(totalDealCashbacks)
	}

	DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
		Where(&DealCashback{UserGUID: userGUID}).Group("shopping_lists.guid").
		Find(&dealCashbacks)

	return dealCashbacks, len(dealCashbacks)
}

// GetByUserGUIDAndDealGUIDGroupByShoppingList function used to retrieve deal cashbacks by user GUID and deal GUID.
func (dcr *DealCashbackRepository) GetByUserGUIDAndDealGUIDGroupByShoppingList(userGUID, dealGUID, pageNumber, pageLimit, relations string) ([]*DealCashback, int) {
	dealCashbacks := []*DealCashback{}

	DB := dcr.DB.Model(&DealCashback{})

	offset := dcr.SetOffsetValue(pageNumber, pageLimit)

	if relations != "" {
		DB = dcr.LoadRelations(DB, relations)
	}

	// This code will be executed when the request comewant API to return result deal cashbacks in paginaton.
	// It will check if page limit parameter not empty.
	if pageLimit != "" {
		totalDealCashbacks := []*DealCashback{}

		// Query to count total number of deal cashback filter by user GUID and deal GUID and group by shopping list guid.
		DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
			Where(&DealCashback{UserGUID: userGUID, DealGUID: dealGUID}).Group("shopping_lists.guid").
			Find(&totalDealCashbacks)

		// Query to retrieve deal cashbacks
		DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
			Where(&DealCashback{UserGUID: userGUID, DealGUID: dealGUID}).Group("shopping_lists.guid").
			Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

		return dealCashbacks, len(totalDealCashbacks)
	}

	// This code will be executed when the request want API to return deal cashbacks without pagination.
	DB.Joins("left join shopping_lists on shopping_lists.guid = deal_cashbacks.shopping_list_guid").
		Where(&DealCashback{UserGUID: userGUID, DealGUID: dealGUID}).Group("shopping_lists.guid").
		Find(&dealCashbacks)

	return dealCashbacks, len(dealCashbacks)
}

// GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID function used to retrieve deal cashback by
// deal cashback transaction GUID and group the result by shopping list.
func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndGroupByShoppingListGUID(dealCashbackTransactionGUID *string) []*DealCashback {
	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID}).Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

// GetByUserGUIDAndTransactionStatusGroupByShoppingListGUID function used to retrieve deal cashback by
// user GUID and tranasaction status and group the result by shopping list.
func (dcr *DealCashbackRepository) GetByUserGUIDAndTransactionStatusGroupByShoppingListGUID(userGUID, transactionStatus string) []*DealCashback {
	dealCashbacks := []*DealCashback{}

	DB := dcr.DB.Model(&DealCashback{}).Where(&DealCashback{UserGUID: userGUID})

	if transactionStatus == "pending" || transactionStatus == "approved" || transactionStatus == "reject" || transactionStatus == "partial_success" {
		DB = DB.Joins("LEFT JOIN deal_cashback_transactions ON deal_cashback_transactions.guid = deal_cashback_transaction_guid").
			Joins("LEFT JOIN transactions ON transactions.GUID = deal_cashback_transactions.transaction_guid").
			Joins("LEFT JOIN transaction_statuses ON transaction_statuses.guid = transactions.transaction_status_guid").
			Where("transaction_statuses.slug = ?", transactionStatus)
	} else if transactionStatus == "notempty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NOT NULL")
	} else if transactionStatus == "empty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NULL")
	}

	DB.Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

// GetByDealCashbackTransactionGUIDAndShoppingListGUID function used to retrieve deal cashback by transaction
// GUID and shopping list GUID.
func (dcr *DealCashbackRepository) GetByDealCashbackTransactionGUIDAndShoppingListGUID(dealCashbackTransactionGUID *string,
	shoppingListGUID, relations string) []*DealCashback {

	dealCashbacks := []*DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealCashbackTransactionGUID: dealCashbackTransactionGUID, ShoppingListGUID: shoppingListGUID}).
		Group("shopping_list_guid").Find(&dealCashbacks)

	return dealCashbacks
}

// GetByUserGUIDAndShoppingListGUIDAndDealGUID function used to retrieve deal cashback by user GUID, shopping list
// GUID and deal GUID.
func (dcr *DealCashbackRepository) GetByUserGUIDAndShoppingListGUIDAndDealGUID(userGUID, shoppingListGUID, dealGUID string) *DealCashback {
	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID, DealGUID: dealGUID}).First(&dealCashback)

	return dealCashback
}

// CountByDealGUIDAndUserGUID function used to count total number of deal cashback by deal GUID and user GUID.
func (dcr *DealCashbackRepository) CountByDealGUIDAndUserGUID(dealGUID, userGUID string) int {
	var totalNumberOfUserDealCashback int

	dcr.DB.Model(&DealCashback{}).Where(DealCashback{DealGUID: dealGUID, UserGUID: userGUID}).Count(&totalNumberOfUserDealCashback)

	return totalNumberOfUserDealCashback
}

// CountByDealGUID function used to count total number of deal cashback by deal GUID.
func (dcr *DealCashbackRepository) CountByDealGUID(dealGUID string) int {
	var totalDealCashback int

	dcr.DB.Model(&DealCashback{}).Where(&DealCashback{DealGUID: dealGUID}).Count(&totalDealCashback)

	return totalDealCashback
}

// CalculateTotalCashbackAmountFromDealCashbackAddedTolist function used to sum all of cashback amount for deal cashback already
// added to list.
func (dcr *DealCashbackRepository) CalculateTotalCashbackAmountFromDealCashbackAddedTolist(userGUID string) float64 {

	type DealCashback struct {
		TotalAmountOfCashback float64 `json:"total_amount_of_cashback"`
	}

	dealCashback := &DealCashback{}

	dcr.DB.Model(&DealCashback{}).Select("sum(ads.cashback_amount) as total_amount_of_cashback").
		Joins("left join ads on ads.guid = deal_cashbacks.deal_guid").Where("user_guid = ?", userGUID).
		Where("deal_cashbacks.deal_cashback_transaction_guid IS NULL").Scan(dealCashback)

	return dealCashback.TotalAmountOfCashback
}

// GetByUserGUIDAndTransactionStatus function used to retrieve deal cashback by user GUID and transaction status slug.
func (dcr *DealCashbackRepository) GetByUserGUIDAndTransactionStatus(userGUID, transactionStatus, pageNumber, pageLimit,
	relations string) ([]*DealCashback, int) {

	dealCashbacks := []*DealCashback{}

	offset := dcr.SetOffsetValue(pageNumber, pageLimit)

	DB := dcr.DB.Model(&DealCashback{})

	if relations != "" {
		DB = dcr.LoadRelations(DB, relations)
	}

	DB = DB.Where(DealCashback{UserGUID: userGUID})

	if transactionStatus == "pending" || transactionStatus == "approved" || transactionStatus == "reject" || transactionStatus == "partial_success" {
		DB = DB.Joins("LEFT JOIN deal_cashback_transactions ON deal_cashback_transactions.guid = deal_cashback_transaction_guid").
			Joins("LEFT JOIN transactions ON transactions.GUID = deal_cashback_transactions.transaction_guid").
			Joins("LEFT JOIN transaction_statuses ON transaction_statuses.guid = transactions.transaction_status_guid").
			Where("transaction_statuses.slug = ?", transactionStatus)
	} else if transactionStatus == "notempty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NOT NULL")
	} else if transactionStatus == "empty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NULL")
	}

	var totalDealCashback *int

	DB.Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

	DB.Count(&totalDealCashback)

	return dealCashbacks, *totalDealCashback
}

// GetByUserGUIDShoppingListGUIDAndTransactionStatus function used to retrieve deal cashback by user GUID, shopping list GUID
// and transaction status slug.
func (dcr *DealCashbackRepository) GetByUserGUIDShoppingListGUIDAndTransactionStatus(userGUID, shoppingListGUID, transactionStatus,
	pageNumber, pageLimit, relations string) ([]*DealCashback, int) {

	dealCashbacks := []*DealCashback{}

	DB := dcr.DB.Model(&DealCashback{})

	if relations != "" {
		DB = dcr.LoadRelations(DB, relations)
	}

	DB = DB.Where(DealCashback{UserGUID: userGUID, ShoppingListGUID: shoppingListGUID})

	if transactionStatus == "pending" || transactionStatus == "approved" || transactionStatus == "reject" || transactionStatus == "partial_success" {
		DB = DB.Joins("LEFT JOIN deal_cashback_transactions ON deal_cashback_transactions.guid = deal_cashback_transaction_guid").
			Joins("LEFT JOIN transactions ON transactions.GUID = deal_cashback_transactions.transaction_guid").
			Joins("LEFT JOIN transaction_statuses ON transaction_statuses.guid = transactions.transaction_status_guid").
			Where("transaction_statuses.slug = ?", transactionStatus)
	} else if transactionStatus == "notempty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NOT NULL")
	} else if transactionStatus == "empty" {
		DB = DB.Where("deal_cashback_transaction_guid IS NULL")
	}

	var totalDealCashback *int

	if pageNumber != "" && pageLimit != "" {
		offset := dcr.SetOffsetValue(pageNumber, pageLimit)

		DB.Offset(offset).Limit(pageLimit).Find(&dealCashbacks)

		DB.Count(&totalDealCashback)

		return dealCashbacks, *totalDealCashback
	}

	DB.Find(&dealCashbacks)

	return dealCashbacks, 0
}
