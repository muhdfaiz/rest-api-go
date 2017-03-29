package v11

// CreateCashoutTransaction is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Cashout Transaction Handler. See `Create` function.
type CreateCashoutTransaction struct {
	Amount                float64 `form:"amount" json:"amount" binding:"required,gt=0"`
	BankAccountHolderName string  `form:"bank_account_name" json:"bank_account_name" binding:"required"`
	BankAccountNumber     string  `form:"bank_account_number" json:"bank_account_number" binding:"required,numeric"`
	BankName              string  `form:"bank_name" json:"bank_name" binding:"required"`
	BankCountry           string  `form:"bank_country" json:"bank_country" binding:"required"`
}
