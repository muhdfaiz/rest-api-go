package v1

import "time"

type EventService struct {
	EventRepository        EventRepositoryInterface
	DealCashbackRepository DealCashbackRepositoryInterface
	DealService            DealServiceInterface
}

func (es *EventService) GetAllIncludingDeals(userGUID string) []*Event {
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

		deal := es.DealService.SetAddTolistInfoAndItemsAndGrocerExclusiveForDeals(filteredDealsQuota, userGUID)

		event.Deals = deal

	}

	return events
}
