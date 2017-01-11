package v1_1

type CreateTransaction struct {
	UserGUID              string  `json:"user_guid"`
	TransactionTypeGUID   string  `json:"transaction_type_guid"`
	TransactionStatusGUID string  `json:"transaction_status_guid"`
	Amount                float64 `json:"amount"`
	ReferenceID           string  `json:"reference_id"`
}
