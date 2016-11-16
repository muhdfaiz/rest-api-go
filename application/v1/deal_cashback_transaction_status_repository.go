package v1

import "github.com/jinzhu/gorm"

type DealCashbackTransactionStatusRepositoryInterface interface {
	GetBySlug(slug string) *DealCashbackTransactionStatus
}

type DealCashbackTransactionStatusRepository struct {
	DB *gorm.DB
}

func (dctsr *DealCashbackTransactionStatusRepository) GetBySlug(slug string) *DealCashbackTransactionStatus {
	dealCashbackTransactionStatus := &DealCashbackTransactionStatus{}

	dctsr.DB.Model(&DealCashbackTransactionStatus{}).Where(&DealCashbackTransactionStatus{Slug: slug}).First(&dealCashbackTransactionStatus)

	return dealCashbackTransactionStatus
}
