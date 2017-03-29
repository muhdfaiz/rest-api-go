package v11

import (
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type GenericResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	LastUpdate     *time.Time  `json:"last_update"`
	Data           interface{} `json:"data"`
}

type GenericTransformer struct{}

func (gt *GenericTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *GenericResponse {
	genericResponse := &GenericResponse{}
	genericResponse.TotalCount = totalData
	genericResponse.Data = data

	limitInt, _ := strconv.Atoi(limit)

	if limitInt != -1 && limitInt != 0 {
		genericResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
	}

	generics := data.([]*Generic)

	if len(generics) == 0 {
		genericResponse.LastUpdate = nil
		return genericResponse
	}

	genericResponse.LastUpdate = &generics[0].UpdatedAt

	return genericResponse
}
