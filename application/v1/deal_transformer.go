package v1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type DealResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	Data           interface{} `json:"data"`
}

type DealTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *DealResponse
}

type DealTransformer struct{}

func (dt *DealTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *DealResponse {
	dealResponse := &DealResponse{}
	dealResponse.TotalCount = totalData
	dealResponse.Data = data

	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)

		if limitInt != -1 {
			dealResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
		}
	}

	return dealResponse
}
