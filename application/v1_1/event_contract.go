package v1_1

// EventServiceInterface is a contract that defines the method needed for Event Service.
type EventServiceInterface interface {
	GetAllIncludingDeals(userGUID string) []*Event
}

// EventRepositoryInterface is a contract that defines the method needed for Event Repository.
type EventRepositoryInterface interface {
	GetAllIncludingRelations(todayDateInGMT8 string) []*Event
}
