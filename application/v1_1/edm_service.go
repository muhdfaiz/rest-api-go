package v1_1

import (
	"strconv"

	"encoding/json"

	"bitbucket.org/cliqers/shoppermate-api/services/email"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type EdmService struct {
	EmailService           email.EmailServiceInterface
	DealService            DealServiceInterface
	FeaturedDealRepository FeaturedDealRepositoryInterface
}

type EdmVariablesStructure struct {
	Name    string      `json:"name"`
	Content interface{} `json:"content"`
}

func (es *EdmService) SendEdmForInsufficientFunds(userGUID string, data SendEdmInsufficientFunds) *systems.ErrorData {
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

	latestDeals, _ := es.DealService.GetAvailableDealsForRegisteredUser(userGUID, "", data.Latitude, data.Longitude, "1", "3", "")

	var latestDealVariables []map[string]string

	for _, latestDeal := range latestDeals {
		latestDealVariable := map[string]string{
			"img":   latestDeal.Img,
			"name":  latestDeal.Name,
			"price": strconv.FormatFloat(latestDeal.CashbackAmount, 'f', 2, 64),
		}

		latestDealVariables = append(latestDealVariables, latestDealVariable)
	}

	variable1 := &EdmVariablesStructure{
		Name:    "user_fullname",
		Content: data.Name,
	}

	variable2 := &EdmVariablesStructure{
		Name:    "products",
		Content: featuredDealVariables,
	}

	variable3 := &EdmVariablesStructure{
		Name:    "deals",
		Content: latestDealVariables,
	}

	edmVariables := make([]interface{}, 3)

	edmVariables[0] = variable1
	edmVariables[1] = variable2
	edmVariables[2] = variable3

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

	if error != nil {
		return error1
	}

	return nil
}
