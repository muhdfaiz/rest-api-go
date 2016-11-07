package v1

import "github.com/jinzhu/gorm"

type DealRepositoryInterface interface {
	GetDealsByCategoryAndValidStartEndDate(todayDateInGMT8 string, shoppingListItem *ShoppingListItem) []*Deal
	GetDealsByValidStartEndDate(todayDateInGMT8 string) []*Deal
	GetDealByGUID(dealGUID string) *Deal
	GetDealByIDWithRelations(dealID int, relations string) *Ads
	GetAllDealsWithinValidRangeAndQuota(latitude float64, longitude float64, offset string, limit string, relations string) ([]*Deal, int)
	GetDealByGUIDAndValidStartEndDate(dealGUID string, todayDateInGMT8 string) *Deal
}

type DealRepository struct {
	DB                    *gorm.DB
	GrocerLocationService GrocerLocationServiceInterface
}

// GetDealsByCategoryAndValidStartEndDate used to retrieve deals by category and between start date & end date
func (dr *DealRepository) GetDealsByCategoryAndValidStartEndDate(todayDateInGMT8 string, shoppingListItem *ShoppingListItem) []*Deal {
	deals := []*Deal{}

	dr.DB.Model(&Deal{}).Where("category = ? AND start_date <= ? AND end_date > ? AND status = ?", shoppingListItem.Category,
		todayDateInGMT8, todayDateInGMT8, "publish").Find(&deals)

	return deals
}

// GetDealsByValidStartEndDate used to retrieve deals that still valid between start date & end date
func (dr *DealRepository) GetDealsByValidStartEndDate(todayDateInGMT8 string) []*Deal {
	deals := []*Deal{}

	dr.DB.Model(&Deal{}).Where("start_date <= ? AND end_date > ? AND status = ?", todayDateInGMT8, todayDateInGMT8, "publish").Find(&deals)

	return deals
}

// GetDealByGUID used to retrieve deal by GUID
func (dr *DealRepository) GetDealByGUID(dealGUID string) *Deal {
	deal := &Deal{}

	dr.DB.Model(&Deal{}).Where("guid = ?", dealGUID).Find(&deal)

	return deal
}

// GetDealByIDWithRelations used to retrieve deal by ID including the relations like grocers, grocer locations and item
// Note: Need to use `Ads` model due to GORM ORM set the column name based on struct name on pivot table
// For example if the model name is `Deal` then GORM ORM will use deal_id to match inside pivot table `ads_grocer`
// Right now, column name inside pivot table `ads_grocer` is `ads_id`. So must use `Ads` model.
func (dr *DealRepository) GetDealByIDWithRelations(dealID int, relations string) *Ads {
	deal := &Ads{}

	DB := dr.DB.Model(&Ads{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Where(&Ads{ID: dealID}).Preload("Grocers", func(db *gorm.DB) *gorm.DB {
		return db.Where("ads_id = ?", dealID).Group("grocer_id")
	}).Find(&deal)

	for key, dealGrocer := range deal.Grocers {
		grocerLocations := []*GrocerLocation{}

		joinQuery := "RIGHT JOIN ads_grocer ON grocer_location.id = ads_grocer.grocer_location_id WHERE ads_grocer.ads_id = ? AND ads_grocer.grocer_id = ?"

		DB.Table("grocer_location").Joins(joinQuery, deal.ID, dealGrocer.ID).Scan(&grocerLocations)

		deal.Grocers[key].GrocerLocations = grocerLocations
	}

	return deal

}

// GetAllDealsWithinValidRangeAndQuota used to retrieve deal within valid range (10KM) and the quota still available
// including the relations like grocers, grocer locations and item
func (dr *DealRepository) GetAllDealsWithinValidRangeAndQuota(latitude float64, longitude float64, pageNumber string, pageLimit string, relations string) ([]*Deal, int) {
	dealsWithin10KM := []*Deal{}

	offset := SetOffsetValue(pageNumber, pageLimit)
	sqlQueryStatement := `SELECT deals.*,
       count(deal_cashbacks.deal_guid) AS total_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.category,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				grocer_location.name AS grocer_location_name,
				grocer_location.lat,
				grocer_location.lng,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS distance
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		WHERE ads.status = "publish"
		GROUP BY ads_guid
		HAVING distance <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT JOIN deal_cashbacks ON deal_cashbacks.deal_guid = ads_guid
	GROUP BY deal_cashbacks.deal_guid
	HAVING total_cashback < ads_quota LIMIT ?
	OFFSET ?;`

	dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, 10, pageLimit, offset).Scan(&dealsWithin10KM)

	type TotalDeal struct {
		Total int `json:"total"`
	}

	total := &TotalDeal{}

	dr.DB.Raw(`SELECT FOUND_ROWS() as total;`).Scan(total)

	return dealsWithin10KM, total.Total
}

// GetDealByGUIDAndValidStartEndDate used to retrieve deal by GUID
func (dr *DealRepository) GetDealByGUIDAndValidStartEndDate(dealGUID string, todayDateInGMT8 string) *Deal {
	deal := &Deal{}

	dr.DB.Model(&Deal{}).Where("guid = ? AND start_date <= ? AND end_date > ?", dealGUID, todayDateInGMT8, todayDateInGMT8).Find(&deal)

	return deal
}
