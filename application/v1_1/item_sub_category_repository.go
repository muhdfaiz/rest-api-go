package v1_1

import (
	"os"

	"github.com/jinzhu/gorm"
)

type ItemSubCategoryRepository struct {
	DB *gorm.DB
}

// GetByID function used to retrieve Item Sub Category by ID
func (iscr *ItemSubCategoryRepository) GetByID(id int) *ItemSubCategory {
	itemSubCategory := &ItemSubCategory{}

	iscr.DB.Model(&ItemSubCategory{}).Where(&ItemSubCategory{ID: id}).First(&itemSubCategory)

	return itemSubCategory
}

// GetByGUID function used to retrieve Item Sub Category by GUID
func (iscr *ItemSubCategoryRepository) GetByGUID(guid string) *ItemSubCategory {
	itemSubCategory := &ItemSubCategory{}

	iscr.DB.Model(&ItemSubCategory{}).Where(&ItemSubCategory{GUID: guid}).First(&itemSubCategory)

	return itemSubCategory
}

// GetSubCategoriesForCategoryThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota function used to retrieve item subcategory for category
// those only have deals within range (10KM radius), date range, number of deal added to list by user not exceed number of deal allow per user
// and quota.
func (iscr *ItemSubCategoryRepository) GetSubCategoriesForCategoryThoseHaveDealsWithinRangeAndDateRangeAndUserLimitAndQuota(userGUID, categoryGUID, currentDateInGMT8 string, latitude,
	longitude float64) []*ItemSubCategory {

	uniqueSubCategories := []*ItemSubCategory{}

	sqlQueryStatement := `SELECT deals.*, count(deal_cashbacks.deal_guid) AS total_deal_cashback,
		(SELECT count(*) FROM deal_cashbacks WHERE deal_cashbacks.deal_guid = deals.ads_guid AND user_guid = ?) AS total_user_deal_cashback
		FROM
		(SELECT subcategory.id,
				subcategory.guid,
				category.name AS category_name,
				subcategory.name,
				ads.id AS ads_id,
				ads.guid AS ads_guid,
				ads.advertiser_id AS ads_advertiser_id,
				ads.campaign_id as ads_campaign_id,
				ads.item_id AS ads_item_id,
				ads.category_id AS ads_category_id,
				ads.img AS ads_img,
				ads.front_name AS ads_front_name,
				ads.name AS ads_name,
				ads.body AS ads_body,
				ads.start_date AS ads_start_date,
				ads.end_date AS ads_end_date,
				ads.positive_tag AS ads_positive_tag,
				ads.negative_tag AS ads_negative_tag,
				ads.time AS ads_time,
				ads.refresh_period AS ads_refresh_period,
				ads.perlimit AS ads_perlimit,
				ads.cashback_amount AS ads_cashback_amount,
				ads.quota AS ads_quota,
				ads.status AS ads_status,
				ads.grocer_exclusive,
				ads.terms,
				ads.created_at as deal_created_time,
				grocer_location.name AS nearest_grocer_name,
				grocer_location.lat AS nearest_grocer_latitude,
				grocer_location.lng AS nearest_grocer_longitude,
				( min(6373 * acos ( cos (radians(?)) * cos(radians(grocer_location.lat)) * cos(radians(grocer_location.lng) - radians(?)) + sin (radians(?)) * sin(radians(grocer_location.lat)) ))) AS distance
		FROM ads
		INNER JOIN ads_grocer ON ads.id = ads_grocer.ads_id
		INNER JOIN grocer_location ON grocer_location.id = ads_grocer.grocer_location_id
		LEFT JOIN item ON item.id = ads.item_id
		LEFT JOIN category ON category.id = item.category_id
		LEFT JOIN subcategory ON subcategory.id = item.subcategory_id
		WHERE ads.status = "publish" AND ads.start_date <= ? AND ads.end_date > ? AND category.guid = ?
		GROUP BY subcategory.name
		HAVING distance <= ?
		ORDER BY ads.created_at DESC) AS deals
	LEFT OUTER JOIN deal_cashbacks ON ads_guid = deal_cashbacks.deal_guid
	GROUP BY ads_guid
	HAVING total_deal_cashback < ads_quota AND total_user_deal_cashback < ads_perlimit
	ORDER BY deal_created_time DESC`

	iscr.DB.Raw(sqlQueryStatement, userGUID, latitude, longitude, latitude, currentDateInGMT8, currentDateInGMT8,
		categoryGUID, os.Getenv("MAX_DEAL_RADIUS_IN_KM")).Scan(&uniqueSubCategories)

	return uniqueSubCategories
}
