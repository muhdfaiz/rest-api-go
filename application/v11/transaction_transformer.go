package v11

import (
	"net/http"
	"strconv"

	"bitbucket.org/cliqers/shoppermate-api/systems"
)

// TransactionResponse define the response structure of transaction data.
// This useful if API want to display paginated response that contain total data and links to current page, last page, first page & next page.
type TransactionResponse struct {
	systems.TotalData
	*systems.Links `json:"links,omitempty"`
	Data           interface{} `json:"data"`
}

// TransactionTransformer is a transformer that used to transform collection of transaction into Transaction Response structure.
type TransactionTransformer struct{}

// transformCollection used to set additional data according to Transaction Response.
func (tt *TransactionTransformer) transformCollection(request *http.Request, data interface{}, totalData int, limit string) *TransactionResponse {
	transactionResponse := &TransactionResponse{}
	transactionResponse.TotalCount = totalData
	transactionResponse.Data = data

	// Convert limit to int if not empty
	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)

		if limitInt != -1 && limitInt != 0 {
			// If limit not euqal to `-1` and `0' set transaction links. Transaction links contain link to current page,
			// first page, next page and last page.
			transactionResponse.Links = PaginationReponse.BuildPaginationLinks(request, totalData)
		}
	}

	return transactionResponse
}
