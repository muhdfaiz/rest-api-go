package v11

import (
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type OccasionResponse struct {
	systems.TotalData
	LastUpdate *time.Time  `json:"last_update"`
	Data       interface{} `json:"data"`
}

type OccasionTransformer struct{}

func (ot *OccasionTransformer) TransformCollection(data interface{}, totalData int) *OccasionResponse {
	occasions := data.([]*Occasion)

	occasionResponse := &OccasionResponse{}
	occasionResponse.TotalCount = totalData
	occasionResponse.Data = data

	if len(occasions) == 0 {
		occasionResponse.LastUpdate = nil
		return occasionResponse
	}

	occasionResponse.LastUpdate = &occasions[0].UpdatedAt

	return occasionResponse
}
