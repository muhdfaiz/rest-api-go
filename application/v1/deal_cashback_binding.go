package v1

// AddDealToShoppingListAndCashback will bind request data based on header content type
type CreateDealCashback struct {
	ShoppingListGUID string `form:"user_guid" json:"shopping_list_guid" binding:"required,uuid5"`
	DealGUID         string `form:"uuid" json:"deal_guid" binding:"required,uuid4"`
}
