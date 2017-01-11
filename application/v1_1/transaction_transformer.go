package v1_1

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

type TransactionResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	Data           interface{} `json:"data"`
}

type TransactionTransformerInterface interface {
	transformCollection(currentURI *http.Request, data interface{}, totalData int, limit string) *TransactionResponse
}

type TransactionTransformer struct{}

func (tt *TransactionTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *TransactionResponse {
	transactionResponse := &TransactionResponse{}
	transactionResponse.TotalCount = totalData
	transactionResponse.Data = data

	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)

		if limitInt != -1 && limitInt != 0 {
			transactionResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
		}
	}

	return transactionResponse
}
