package v11

// CreateDealCashback is a request data binding that will be used to bind request body to struct.
// When API receives request with header `application/json`, GIN will used `json` tag to find the data.
// When API receives request with header `multipart/form-data` or `application/x-www-form-urlencoded`,
// GIN will used `form` tag to find the data.
// Used in Deal Cashback Handler. See `Create` function.
type CreateDealCashback struct {
	ShoppingListGUID string `form:"shopping_list_guid" json:"shopping_list_guid" binding:"required,uuid5"`
	DealGUID         string `form:"deal_guid" json:"deal_guid" binding:"required"`
}
