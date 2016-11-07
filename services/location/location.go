package location

import "math"

type LocationServiceInterface interface {
	CalculateDistanceInMeter(lat1, lon1, lat2, lon2 float64) float64
	CalculateDistanceInKilometer(latitude1, longitude1, latitude2, longitude2 float64) float64
}

type LocationService struct{}

// CalculateDistanceInMeter function returns the distance (in meters) between two points of
// a given longitude and latitude relatively accurately (using a spherical
// approximation of the Earth) through the Haversin Distance Formula for
// great arc distance on a sphere with accuracy for small distances.
//
// Point coordinates are supplied in degrees and converted into rad. in the func
//
// Distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func (ls *LocationService) CalculateDistanceInMeter(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var lat1, lon1, lat2, lon2, r float64
	lat1 = latitude1 * math.Pi / 180
	lon1 = longitude1 * math.Pi / 180
	lat2 = latitude2 * math.Pi / 180
	lon2 = longitude2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := ls.hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*ls.hsin(lon2-lon1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func (ls *LocationService) CalculateDistanceInKilometer(latitude1, longitude1, latitude2, longitude2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var lat1, lon1, lat2, lon2, r float64
	lat1 = latitude1 * math.Pi / 180
	lon1 = longitude1 * math.Pi / 180
	lat2 = latitude2 * math.Pi / 180
	lon2 = longitude2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := ls.hsin(lat2-lat1) + math.Cos(lat1)*math.Cos(lat2)*ls.hsin(lon2-lon1)

	return 2 * r * math.Asin(math.Sqrt(h)) / 1000
}

// haversin(θ) function
func (ls *LocationService) hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
