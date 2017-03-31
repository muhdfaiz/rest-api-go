package v11

import (
	"os"

	"github.com/jinzhu/gorm"
)

// DealRepository will handle all CRUD function for task related to Deal resource.
type DealRepository struct {
	BaseRepository
	DB                    *gorm.DB
	GrocerLocationService GrocerLocationServiceInterface
}

// SumCashbackAmount function used to sum total amount of deal cashback from multiple deal GUID.
func (dr *DealRepository) SumCashbackAmount(dealGUIDs []string) float64 {
	type Ads struct {
		TotalCashbackAmount float64 `json:"total_cashback_amount"`
	}

	deal := &Ads{}

	dr.DB.Model(&Ads{}).Select("sum(cashback_amount) as total_cashback_amount").Where("guid in (?)", dealGUIDs).Find(&deal)

	return deal.TotalCashbackAmount
}

// GetDealByGUID used to retrieve deal details from database using deal GUID.
func (dr *DealRepository) GetDealByGUID(dealGUID string) *Deal {
	deal := &Deal{}

	dr.DB.Model(&Deal{}).Where("guid = ?", dealGUID).Find(&deal)

	return deal
}

// GetDealByGUIDAndStartEndDateAndPublishStatus used to retrieve deal by GUID and the deal start date and end date within the current date.
func (dr *DealRepository) GetDealByGUIDAndStartEndDateAndPublishStatus(dealGUID, todayDateInGMT8 string) *Deal {
	deal := &Deal{}

	dr.DB.Model(&Deal{}).Where("guid = ? AND start_date <= ? AND end_date > ? AND status = ?", dealGUID, todayDateInGMT8, todayDateInGMT8, "publish").
		Find(&deal)

	return deal
}

// GetDealByIDWithRelations used to retrieve deal by ID including the relations like grocers, grocer locations and item
// Note: Need to use `Ads` model because GORM ORM set the column name based on struct name on pivot table
// For example if the model name is `Deal` then GORM ORM will use deal_id to match inside pivot table `ads_grocer`
// Right now, column name inside pivot table `ads_grocer` is `ads_id`. So must use `Ads` model.
func (dr *DealRepository) GetDealByIDWithRelations(dealID int, relations string) *Ads {
	deal := &Ads{}

	DB := dr.DB.Model(&Ads{})

	if relations != "" {
		DB = dr.LoadRelations(DB, relations)
	}

	DB.Where(&Ads{ID: dealID}).Preload("Grocers", func(db *gorm.DB) *gorm.DB {
		return db.Where("ads_grocer.ads_id = ? AND grocer.status = ?", dealID, "publish").Group("ads_grocer.grocer_id")
	}).Find(&deal)

	for key, dealGrocer := range deal.Grocers {
		grocerLocations := []*GrocerLocation{}

		joinQuery := "RIGHT JOIN ads_grocer ON grocer_location.id = ads_grocer.grocer_location_id WHERE ads_grocer.ads_id = ? AND ads_grocer.grocer_id = ? AND ads_grocer.grocer_id"

		DB.Table("grocer_location").Joins(joinQuery, deal.ID, dealGrocer.ID).Scan(&grocerLocations)

		deal.Grocers[key].GrocerLocations = grocerLocations
	}

	return deal

}

// GetDealByGUIDAndUserGUIDWithinDateRangeAndValidQuotaAndLimitPerUserAndPublished function used to retrieve deal by GUID with conditions below.
// Deal is valid if the deal start date and end date within current date.
// Deal is valid if the deal status is published.
// Deal is valid if the deal quota still not reach the limit.
// Deal is valid if the user still not reach deal limit per user.
func (dr *DealRepository) GetDealByGUIDAndUserGUIDWithinDateRangeAndValidQuotaAndLimitPerUserAndPublished(userGUID, dealGUID, currentDateInGMT8 string) *Deal {
	deal := &Deal{}

	sqlQueryStatement := `SELECT deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	   (SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at
		FROM ads
		WHERE ads.status = "publish" AND ads.guid = ? AND ads.start_date <= ? AND ads.end_date > ?
		GROUP BY ads.guid
		ORDER BY ads.created_at DESC) AS deals
	LEFT JOIN deal_cashbacks ON deal_cashbacks.deal_guid = ads_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	dr.DB.Raw(sqlQueryStatement, userGUID, dealGUID, currentDateInGMT8, currentDateInGMT8).Scan(deal)

	return deal
}

// GetTodayDealsWithPublishStatusByShoppingListItemCategoryName used to retrieve available deals related to shopping list item.
// Use this function when you want to display related deals in shopping list items.
//
// Find deals based on criteria below:
// - Shopping List Item category name in parameter must match with deal category name.
// - Today date is within deal start date and end date.
// - Deal status must be publish.
func (dr *DealRepository) GetTodayDealsWithPublishStatusByShoppingListItemCategoryName(todayDateInGMT8 string, shoppingListItem *ShoppingListItem) []*Deal {
	deals := []*Deal{}

	dr.DB.Model(&Deal{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Where(&ItemCategory{Name: shoppingListItem.Category})
	}).Where("start_date <= ? AND end_date > ? AND status = ?", todayDateInGMT8, todayDateInGMT8, "publish").Find(&deals)

	return deals
}

// GetTodayDealsWithPublishStatus used to retrieve available deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
func (dr *DealRepository) GetTodayDealsWithPublishStatus(todayDateInGMT8 string) []*Deal {
	deals := []*Deal{}

	dr.DB.Model(&Deal{}).Where("start_date <= ? AND end_date > ? AND status = ?", todayDateInGMT8, todayDateInGMT8, "publish").Find(&deals)

	return deals
}

// GetTodayDealsWithValidQuotaAndNearUserLocation used to retrieve available deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetTodayDealsWithValidQuotaAndNearUserLocation(currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*,
       count(deal_cashbacks.deal_guid) AS total_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ?
		GROUP BY ads.guid
		ORDER BY ads.created_at DESC) AS deals
	LEFT JOIN deal_cashbacks ON deal_cashbacks.deal_guid = ads_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8, pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsWithinRangeAndDateRangeAndQuota used to retrieve available deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsWithinRangeAndDateRangeAndQuota(latitude, longitude float64, currentDateInGMT8,
	pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*,
       count(deal_cashbacks.deal_guid) AS total_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6371 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ?
		GROUP BY ads.guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT JOIN deal_cashbacks ON deal_cashbacks.deal_guid = ads_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsWithinRangeAndDateRangeAndUserLimitAndQuotaAndName used to retrieve available deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsWithinRangeAndDateRangeAndUserLimitAndQuotaAndName(userGUID, name string, latitude, longitude float64,
	currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN category ON category.id = ads.category_id
		WHERE ads.status = "publish" AND ads.name LIKE ? AND ads.start_date <= ? AND ads.end_date > ?
		GROUP BY ads_guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, "%"+name+"%", currentDateInGMT8,
			currentDateInGMT8, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, "%"+name+"%", currentDateInGMT8,
		currentDateInGMT8, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, "%"+name+"%", currentDateInGMT8,
		currentDateInGMT8, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsForCategoryWithinDateRangeAndQuota used to retrieve available deals for specific category using category name.
// For example in Shoppermate has `Drinks` category. Use this function to retrieve deals available in `Drinks` category.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal category name must match with category name in the parameter.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsForCategoryWithinDateRangeAndQuota(category, currentDateInGMT8,
	pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*,
       count(deal_cashbacks.deal_guid) AS total_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at
		FROM ads
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN category ON category.id = item.category_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND category.name = ?
		GROUP BY ads.guid
		ORDER BY ads.created_at DESC) AS deals
	LEFT JOIN deal_cashbacks ON deal_cashbacks.deal_guid = ads_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8, category).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8, category).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, currentDateInGMT8, currentDateInGMT8, category, pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsForCategoryWithinRangeAndDateRangeAndQuota used to retrieve available deals for specific category using category name.
// For example in Shoppermate has `Drinks` category. Use this function to retrieve deals available in `Drinks` category.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal category name must match with category name in the parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsForCategoryWithinRangeAndDateRangeAndQuota(category string, latitude, longitude float64,
	currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback
		FROM
			(SELECT ads.id AS ads_id,
					ads.guid AS ads_guid,
					ads.id,
					ads.guid,
					ads.advertiser_id,
					ads.campaign_id,
					ads.item_id,
					ads.category_id,
					ads.img,
					ads.front_name,
					ads.name,
					ads.body,
					ads.start_date,
					ads.end_date,
					ads.positive_tag,
					ads.negative_tag,
					ads.time,
					ads.refresh_period,
					ads.perlimit,
					ads.perlimit AS ads_perlimit,
					ads.cashback_amount,
					ads.quota AS ads_quota,
					ads.quota,
					ads.status,
					ads.grocer_exclusive,
					ads.terms,
					ads.created_at as deal_created_time,
					ads.created_at,
					ads.updated_at,
					ads.deleted_at,
					category.name AS category_name,
					grocer_location.name AS nearest_grocer_name,
					grocer_location.lat AS nearest_grocer_latitude,
					grocer_location.lng AS nearest_grocer_longitude,
					( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
			FROM ads
			INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
			INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
			LEFT JOIN item ON item.id = ads.item_id
			LEFT JOIN category ON category.id = item.category_id
			WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND category.name = ?
			GROUP BY ads_guid
			HAVING nearest_grocer_distance_in_km <= ?
			ORDER BY ads.created_at DESC) AS deals
		LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
		GROUP BY ads_guid
		HAVING total_deal_cashback < ads_quota
		ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			category, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		category, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		category, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsByCategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota used to retrieve available deals for specific category using category name.
// For example in Shoppermate has `Drinks` category. Use this function to retrieve deals available in `Drinks` category.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal category name must match with category name in the parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsByCategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, category string, latitude, longitude float64,
	currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				category.name AS category_name,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN category ON category.id = item.category_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND category.name = ?
		GROUP BY ads_guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			category, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		category, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		category, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuotaAndCategory used to retrieve multiple deals for specific grocer and category.
// For example: you want retrieve deals available at Tesco Hypermarket and to retrieve deals from category `drink` only.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal grocer ID must match with grocer ID in parameter.
// - Deal category GUID must match with category GUID in parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuotaAndCategory(userGUID, categoryGUID string, grocerID int,
	latitude, longitude float64, currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
		(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
		FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				category.name AS category_name,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN category ON category.id = item.category_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND ads_grocer.grocer_id = ? AND category.guid = ?
		GROUP BY ads_guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			grocerID, categoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		grocerID, categoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		grocerID, categoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// CountDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuota used to count deals available for specific grocer using grocer ID.
// For example you want to count available deals at Tesco Hypermarket.
//
// Count deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal grocer ID must match with grocer ID in parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) CountDealsForGrocerWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID string, grocerID int,
	latitude, longitude float64, currentDateInGMT8 string) int {

	deals := []*Deal{}

	sqlQueryStatement := `SELECT deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				category.name AS category_name,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN category ON category.id = item.category_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND ads_grocer.grocer_id = ?
		GROUP BY ads_guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		grocerID, 10).Scan(&deals)

	return len(deals)
}

// GetDealBySubCategoryGUIDWithinRangeAndDateRangeAndUserLimitAndQuota used to retrieve available deals for specific subcategory using subcategory GUID.
// For example in shoppermate has `Drinks` category and inside that category contains multiple subcategory. Example of subcategory is
// `Carbonated Drinks`, `Coffee` and `Milk`. Use this function to retrieve deals available in subcategory like `Carbonated Drinks`.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal Subcategory GUID must match with subcategory GUID in parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealBySubCategoryGUIDWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, subCategoryGUID string,
	latitude, longitude float64, currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT subcategory.id as subcategory_id,
				subcategory.guid as subcategory_guid,
				subcategory.name as subcategory_name, 
				ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN subcategory ON subcategory.id = item.subcategory_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND subcategory.guid = ?
		GROUP BY ads_id
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			subCategoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		subCategoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		subCategoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}

// GetDealsBySubcategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota used to retrieve available deals for specific subcategory using subcategory Name.
// For example in shoppermate has `Drinks` category and inside that category contains multiple subcategory. Example of subcategory is
// `Carbonated Drinks`, `Coffee` and `Milk`. Use this function to retrieve deals available in subcategory like `Carbonated Drinks`.
//
// Find deals based on criteria below:
// - Deal status must be publish.
// - Today date is within deal start date and end date.
// - Deal subcategory name must match with subcategory name in the parameter.
// - User location must be within valid range (10KM radius).
// - Total number of deal added to list by user not more than deal perlimit.
// - Total number of deal added to list by all user not more than deal quota.
func (dr *DealRepository) GetDealsBySubcategoryNameWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, subcategory string, latitude, longitude float64,
	currentDateInGMT8, pageNumber, pageLimit, relations string) ([]*Deal, int) {

	deals := []*Deal{}

	offset := dr.SetOffsetValue(pageNumber, pageLimit)

	sqlQueryStatement := `SELECT SQL_CALC_FOUND_ROWS deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
	(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
	FROM
		(SELECT ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.id,
				ads.guid,
				ads.advertiser_id,
				ads.campaign_id,
				ads.item_id,
				ads.category_id,
				ads.img,
				ads.front_name,
				ads.name,
				ads.body,
				ads.start_date,
				ads.end_date,
				ads.positive_tag,
				ads.negative_tag,
				ads.time,
				ads.refresh_period,
				ads.perlimit,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount,
				ads.quota AS ads_quota,
				ads.quota,
				ads.status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				ads.created_at,
				ads.updated_at,
				ads.deleted_at,
				subcategory.name AS subcategory_name,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS nearest_grocer_distance_in_km
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN subcategory ON subcategory.id = item.subcategory_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND subcategory.name = ?
		GROUP BY ads_guid
		HAVING nearest_grocer_distance_in_km <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	if pageLimit == "" {
		dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
			subcategory, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&deals)

		return deals, len(deals)
	}

	totalDeal := []*Deal{}

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		subcategory, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&totalDeal)

	sqlQueryStatement = sqlQueryStatement + " LIMIT ? OFFSET ?"

	dr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		subcategory, os.Getenv("MAX_DEAL_RADIUS_IN_KM"), pageLimit, offset).Scan(&deals)

	return deals, len(totalDeal)
}
