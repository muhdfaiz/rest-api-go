package v1_1

// AddDealToShoppingListAndCashback will bind request data based on header content type
type CreateDealCashback struct {
	ShoppingListGUID string `form:"shopping_list_guid" json:"shopping_list_guid" binding:"required,uuid5"`
	DealGUID         string `form:"deal_guid" json:"deal_guid" binding:"required"`
}
