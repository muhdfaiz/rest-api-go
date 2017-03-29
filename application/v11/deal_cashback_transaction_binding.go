package v11

// CreateDealCashbackTransaction is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Deal Cashback Transaction Handler. See `Create` function.
type CreateDealCashbackTransaction struct {
	DealCashbackGuids string `form:"deal_cashback_guids" json:"deal_cashback_guids" binding:"required"`
}
