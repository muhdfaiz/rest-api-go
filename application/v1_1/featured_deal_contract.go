package v1_1

type FeaturedDealRepositoryInterface interface {
	GetActiveFeaturedDeals(pageNumber, pageLimit, relations string) []*Deal
}
