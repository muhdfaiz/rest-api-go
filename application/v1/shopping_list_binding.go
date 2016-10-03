package v1

// CreateShoppingList will bind request data based on header content type
type CreateShoppingList struct {
	OccasionGUID string `form:"occasion_guid" json:"occasion_guid" binding:"required,uuid5"`
	Name         string `form:"name" json:"name" binding:"required"`
}

type UpdateShoppingList struct {
	OccasionGUID string `form:"occasion_guid" json:"occasion_guid" binding:"omitempty,uuid5"`
	Name         string `form:"name" json:"name" binding:"omitempty"`
}
