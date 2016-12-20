package v1

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type GrocerResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	LastUpdate     *time.Time  `json:"last_update"`
	Data           interface{} `json:"data"`
}

type GrocerTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *GrocerResponse
}

type GrocerTransformer struct{}

func (gt *GrocerTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *GrocerResponse {
	items := data.([]*Item)

	grocerResponse := &GrocerResponse{}
	grocerResponse.TotalCount = totalData
	grocerResponse.Data = data

	limitInt, _ := strconv.Atoi(limit)

	if limitInt != -1 && limitInt != 0 {
		grocerResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
	}

	grocerResponse.LastUpdate = &items[0].UpdatedAt

	return grocerResponse
}
