package v11

type FeaturedDealRepositoryInterface interface {
	GetActiveFeaturedDeals(pageNumber, pageLimit, relations string) []*Deal
}
