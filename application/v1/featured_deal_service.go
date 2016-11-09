package v1

import "time"

type EventServiceInterface interface {
	GetAllIncludingValidDeals(userGUID string) []*Event
}

type EventService struct {
	EventRepository        EventRepositoryInterface
	DealCashbackRepository DealCashbackRepositoryInterface
}

func (es *EventService) GetAllIncludingValidDeals(userGUID string) []*Event {
	currentDateInGMT8 := time.Now().UTC().Add(time.Hour * 8).Format("2006-01-02")

	events := es.EventRepository.GetAllIncludingRelations(currentDateInGMT8)

	if len(events) < 1 {
		return nil
	}

	for _, event := range events {
		// Check if total number user add to list more than quota
		filteredDealsQuota := []*Deal{}

		for _, deal := range event.Deals {
			totalNumberOfDealCashback := es.DealCashbackRepository.CountByDealGUID(deal.GUID)

			if totalNumberOfDealCashback < deal.Quota {
				filteredDealsQuota = append(filteredDealsQuota, deal)
			}
		}

		filteredDealsUserLimit := []*Deal{}

		for _, deal := range filteredDealsQuota {
			totalNumberOfUserCashback := es.DealCashbackRepository.CountByDealGUIDAndUserGUID(deal.GUID, userGUID)

			if totalNumberOfUserCashback < deal.Perlimit {
				filteredDealsUserLimit = append(filteredDealsUserLimit, deal)
			}
		}

		event.Deals = filteredDealsUserLimit

	}

	return events
}
