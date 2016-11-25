package v1

type CreateTransaction struct {
	UserGUID            string  `json:"user_guid"`
	TransactionTypeGUID string  `json:"transaction_type_guid"`
	Amount              float32 `json:"amount"`
	ReferenceID         string  `json:"reference_id"`
}
