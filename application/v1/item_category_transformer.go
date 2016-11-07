package v1

import "bitbucket.org/cliqers/shoppermate-api/systems"

type ItemCategoryResponse struct {
	systems.TotalData
	Data interface{} `json:"data"`
}

type ItemCategoryTransformerInterface interface {
	TransformCollection(data interface{}, totalData int) *ItemCategoryResponse
}

type ItemCategoryTransformer struct{}

func (ict *ItemCategoryTransformer) TransformCollection(data interface{}, totalData int) *ItemCategoryResponse {
	itemCategoryResponse := &ItemCategoryResponse{}
	itemCategoryResponse.TotalCount = totalData
	itemCategoryResponse.Data = data

	return itemCategoryResponse
}
