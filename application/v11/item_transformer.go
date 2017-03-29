package v11

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type ItemResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	LastUpdate     *time.Time  `json:"last_update"`
	Data           interface{} `json:"data"`
}

type ItemTransformer struct{}

func (it *ItemTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *ItemResponse {
	items := data.([]*Item)

	itemResponse := &ItemResponse{}
	itemResponse.TotalCount = totalData
	itemResponse.Data = data

	limitInt, _ := strconv.Atoi(limit)

	if limitInt != -1 && limitInt != 0 {
		itemResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
	}

	if len(items) == 0 {
		itemResponse.LastUpdate = nil
		return itemResponse
	}

	itemResponse.LastUpdate = &items[0].UpdatedAt
	return itemResponse

}
