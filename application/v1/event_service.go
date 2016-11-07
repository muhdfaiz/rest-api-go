package v1

type EventServiceInterface interface {
	GetAllIncudingDeals() []*Event
}

type EventService struct {
	EventRepository EventRepositoryInterface
}

func (es *EventService) GetAllIncudingDeals() []*Event {
	deals := es.EventRepository.GeAllIncludingDeals()

	return deals
}
