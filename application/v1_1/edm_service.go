package v1_1

import (
	"strconv"

	"github.com/jinzhu/gorm"

	"encoding/json"

	"bitbucket.org/cliqers/shoppermate-api/services/email"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type EdmService struct {
	EmailService           email.EmailServiceInterface
	DealService            DealServiceInterface
	FeaturedDealRepository FeaturedDealRepositoryInterface
	EdmHistoryRepository   EdmHistoryRepositoryInterface
}

type EdmVariablesStructure struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

func (es *EdmService) SendEdmForInsufficientFunds(dbTransaction *gorm.DB, userGUID string, data SendEdmInsufficientFunds) *systems.ErrorData {
	previousEDMHistory := es.EdmHistoryRepository.GetByUserGUIDAndEventAndCreatedAt(userGUID, "insufficient_funds")

	if previousEDMHistory.UserGUID != "" {
		return Error.GenericError("422", systems.ReachLimitSendEDMInsufficientFundForToday, "Failed To Send Edm For Insufficient Funds.",
			"", "System only allowed send EDM Insufficient Funds one time only per day.")
	}

	featuredDealVariables := es.setFeaturedDealsVariable()

	latestDealVariables := es.setLatestDealsVariable(userGUID, data.Latitude, data.Longitude)

	edmVariables := make([]interface{}, 3)

	variable1 := &EdmVariablesStructure{
		Name:    "user_fullname",
		Content: data.Name,
	}

	edmVariables[0] = variable1

	variable2 := &EdmVariablesStructure{
		Name:    "products",
		Content: featuredDealVariables,
	}

	edmVariables[1] = variable2

	if len(latestDealVariables) > 0 {
		variable3 := &EdmVariablesStructure{
			Name:    "deals",
			Content: latestDealVariables,
		}

		edmVariables[2] = variable3
	}

	jsonString, error := json.Marshal(edmVariables)

	if error != nil {
		return Error.InternalServerError(error, systems.JSONNotValid)
	}

	error1 := es.EmailService.SendTemplate(map[string]string{
		"name":      data.Name,
		"email":     data.Email,
		"template":  "14-shoppermate-insufficient-funds",
		"variables": string(jsonString),
	})

	if error1 != nil {
		return error1
	}

	edmHistory := make(map[string]string)
	edmHistory["guid"] = Helper.GenerateUUID()
	edmHistory["user_guid"] = userGUID
	edmHistory["event"] = "insufficient_funds"

	_, error1 = es.EdmHistoryRepository.Create(dbTransaction, edmHistory)

	if error1 != nil {
		return error1
	}

	return nil
}

func (es *EdmService) setFeaturedDealsVariable() []map[string]string {
	featuredDeals := es.FeaturedDealRepository.GetActiveFeaturedDeals("1", "3", "")

	var featuredDealVariables []map[string]string

	for _, featuredDeal := range featuredDeals {
		featuredDealVariable := map[string]string{
			"img":   featuredDeal.Img,
			"name":  featuredDeal.Name,
			"price": strconv.FormatFloat(featuredDeal.CashbackAmount, 'f', 2, 64),
		}

		featuredDealVariables = append(featuredDealVariables, featuredDealVariable)
	}

	return featuredDealVariables
}

func (es *EdmService) setLatestDealsVariable(userGUID, latitude, longitude string) []map[string]string {
	var latestDealVariables []map[string]string

	if latitude != "" && longitude != "" {
		latestDeals, _ := es.DealService.GetAvailableDealsForRegisteredUser(userGUID, "", latitude, longitude, "1", "3", "")

		for _, latestDeal := range latestDeals {
			latestDealVariable := map[string]string{
				"img":   latestDeal.Img,
				"name":  latestDeal.Name,
				"price": strconv.FormatFloat(latestDeal.CashbackAmount, 'f', 2, 64),
			}

			latestDealVariables = append(latestDealVariables, latestDealVariable)
		}
	}

	return latestDealVariables
}
