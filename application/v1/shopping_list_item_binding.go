package v1

type CreateShoppingListItem struct {
	UserGUID         string `form:"user_guid" json:"user_guid" binding:"omitempty"`
	ShoppingListGUID string `form:"shopping_list_guid" json:"shopping_list_guid" binding:"omitempty,uuid5"`
	Name             string `form:"name" json:"name" binding:"required"`
	Quantity         string `form:"quantity" json:"quantity" binding:"required,numeric,max=3"`
	Remark           string `form:"remark" json:"remark" binding:"omitempty"`
}

type UpdateShoppingListItem struct {
	ShoppingListGUID string `json:"shopping_list_guid" binding:"omitempty,uuid5"`
	Name             string `form:"name" json:"name" binding:"omitempty"`
	Quantity         string `form:"quantity" json:"quantity" binding:"omitempty,numeric,max=3"`
	Remark           string `form:"remark" json:"remark" binding:"omitempty"`
}
