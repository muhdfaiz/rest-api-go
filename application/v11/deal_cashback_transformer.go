package v11

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealCashbackResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	Data           interface{} `json:"data"`
}

type DealCashbackTransformer struct{}

func (dt *DealCashbackTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *DealCashbackResponse {
	dealCashbackResponse := &DealCashbackResponse{}
	dealCashbackResponse.TotalCount = totalData
	dealCashbackResponse.Data = data

	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)

		if limitInt != -1 && limitInt != 0 {
			dealCashbackResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
		}
	}

	return dealCashbackResponse
}
