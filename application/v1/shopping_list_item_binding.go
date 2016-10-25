package v1

type CreateShoppingListItem struct {
	UserGUID         string `form:"user_guid" json:"user_guid" binding:"omitempty"`
	ShoppingListGUID string `form:"shopping_list_guid" json:"shopping_list_guid" binding:"omitempty,uuid5"`
	Name             string `form:"name" json:"name" binding:"required"`
	Quantity         int    `form:"quantity" json:"quantity" binding:"required,gte=0,lte=999"`
	Remark           string `form:"remark" json:"remark" binding:"omitempty"`
}

type UpdateShoppingListItem struct {
	ShoppingListGUID string `form:"shopping_list_guid" json:"shopping_list_guid" binding:"omitempty,uuid5"`
	Name             string `form:"name" json:"name" binding:"omitempty"`
	Quantity         int    `form:"quantity" json:"quantity" binding:"omitempty,gte=0,lte=999"`
	Remark           string `form:"remark" json:"remark" binding:"omitempty"`
	AddedToCart      int    `form:"added_to_cart" json:"added_to_cart" binding:"omitempty,gte=0,lte=1"`
}
