package v1

type CreateDealCashbackTransaction struct {
	DealCashbackGuids string `form:"deal_cashback_guids" json:"deal_cashback_guids" binding:"required"`
}
