package v11

// CreateTransaction is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Deal Cashback Transaction Service. See `CreateTransaction` function.
// Used in Transaction Service. See `CreateTransaction` function.
type CreateTransaction struct {
	UserGUID              string  `json:"user_guid"`
	TransactionTypeGUID   string  `json:"transaction_type_guid"`
	TransactionStatusGUID string  `json:"transaction_status_guid"`
	Amount                float64 `json:"amount"`
	ReferenceID           string  `json:"reference_id"`
}
