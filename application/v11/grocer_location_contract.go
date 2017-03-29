package v11

// GrocerLocationServiceInterface is a contract that defines the method needed for Grocer Location Service.
type GrocerLocationServiceInterface interface {
	GetGrocersFromConvertionLocation(convertionLocation string) []*GrocerLocation
}

// GrocerLocationRepositoryInterface is a contract that defines the method needed for Grocer Location Repository.
type GrocerLocationRepositoryInterface interface {
	GetByIDLatitudeAndLongitude(id int, latitude float64, longitude float64) *GrocerLocation
}
