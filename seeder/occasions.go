package main

import (
	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/systems"
)

func main() {
	Helper := &systems.Helpers{}
	DB := &systems.Database{}
	db := DB.Connect()

	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "gathering", Name: "Gathering", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/gathering.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "travel", Name: "Travel", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/travel.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "jungle_tracking", Name: "Jungle Tracking", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/jungle_tracking.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "picnic", Name: "Picnic", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/picnic.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "field_trip", Name: "Field Trip", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/field_trip.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "household", Name: "Household", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/household.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "cooking", Name: "Cooking", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/cooking.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "outing", Name: "Family Outing", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/family_outing.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "birthday", Name: "Birthday", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/birthday.jpg"})
	db.Create(&v1.Occasion{GUID: Helper.GenerateUUID(), Slug: "festive", Name: "Festive", Image: "https://s3-ap-southeast-1.amazonaws.com/shoppermate/occasion_images/festive.jpg"})
}
