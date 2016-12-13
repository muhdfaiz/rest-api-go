package v1

type CreateCashoutTransaction struct {
	Amount                float64 `form:"amount" json:"amount" binding:"required,gt=0"`
	BankAccountHolderName string  `form:"bank_account_name" json:"bank_account_name" binding:"required"`
	BankAccountNumber     string  `form:"bank_account_number" json:"bank_account_number" binding:"required,numeric"`
	BankName              string  `form:"bank_name" json:"bank_name" binding:"required"`
	BankCountry           string  `form:"bank_country" json:"bank_country" binding:"required"`
}
