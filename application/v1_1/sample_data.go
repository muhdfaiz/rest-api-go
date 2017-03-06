package v1_1

import (
	"time"

	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// SampleData will create all data needed during testing.
type SampleData struct {
	DB *gorm.DB
}

// Settings function used to create sample settings for test database.
func (sd *SampleData) Settings(referralActive, pricePerReferral, maxReferralPerUser string) []*Setting {
	referralActiveSetting := &Setting{
		GUID:  Helper.GenerateUUID(),
		Name:  "Referral Active",
		Slug:  "referral_active",
		Value: referralActive,
	}

	result1 := sd.DB.Create(referralActiveSetting)

	pricePerReferralSetting := &Setting{
		GUID:  Helper.GenerateUUID(),
		Name:  "Price Per Referral",
		Slug:  "referral_price",
		Value: pricePerReferral,
	}

	result2 := sd.DB.Create(pricePerReferralSetting)

	maxReferralPerUserSetting := &Setting{
		GUID:  Helper.GenerateUUID(),
		Name:  "Referral Active",
		Slug:  "max_referral_user",
		Value: maxReferralPerUser,
	}

	result3 := sd.DB.Create(maxReferralPerUserSetting)

	data := make([]*Setting, 3)
	data[0] = result1.Value.(*Setting)
	data[1] = result2.Value.(*Setting)
	data[1] = result3.Value.(*Setting)

	return data
}

// DeviceWithoutUserGUID function used to create sample device without
// user GUID.
func (sd *SampleData) DeviceWithoutUserGUID() (*Device, *systems.ErrorData) {

	deviceGUID := Helper.GenerateUUID()

	device := Device{
		GUID:       deviceGUID,
		UUID:       "F4EE037F5340586BE625EE7678797DD0",
		Os:         "Android",
		Model:      "Xiaomi MI 3W",
		PushToken:  "cnFBoDLh0W4:APA91bHQLLmjPt24uHFzFNNDg8PgrYodmcNJTqTNv2uWuSpsWnAYd8KHhL0layfGOH3I4mtf46cZzymmw_flehg0CgKNa8vMj-D-vGEPTmL_A71_30pC3OZOanXcc3zK7B7x-9-_A9Lc",
		AppVersion: "0.9.9.3-STAGING-BETA",
	}

	result := sd.DB.Create(&device)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*Device), nil
}

// DevicesWithoutUserGUID function used to create sample device with empty user GUID.
func (sd *SampleData) DevicesWithoutUserGUID() ([]*Device, *systems.ErrorData) {
	deviceCollection := make([]*Device, 6)

	device1 := Device{
		GUID:       Helper.GenerateUUID(),
		UUID:       "F4EE037F5340586BE625EE7678797DD0",
		Os:         "Android",
		Model:      "Xiaomi MI 3W",
		PushToken:  "cnFBoDLh0W4:APA91bHQLLmjPt24uHFzFNNDg8PgrYodmcNJTqTNv2uWuSpsWnAYd8KHhL0layfGOH3I4mtf46cZzymmw_flehg0CgKNa8vMj-D-vGEPTmL_A71_30pC3OZOanXcc3zK7B7x-9-_A9Lc",
		AppVersion: "0.9.9.3-STAGING-BETA",
	}

	result1 := sd.DB.Create(&device1)

	if result1.Error != nil || result1.RowsAffected == 0 {
		return nil, Error.InternalServerError(result1.Error, systems.DatabaseError)
	}

	deviceCollection[0] = result1.Value.(*Device)

	device2 := Device{
		GUID:       Helper.GenerateUUID(),
		UUID:       "72529e8b36ef4fa6a58aedfeca821309",
		Os:         "iOS",
		Model:      "iPhone",
		PushToken:  "96f74bf8c22984fc7f1ab4dd269e5bc112f717a13d00fda5de1de57c89702476",
		AppVersion: "0.9.9.3-STAGING-BETA",
	}

	result2 := sd.DB.Create(&device2)

	if result2.Error != nil || result2.RowsAffected == 0 {
		return nil, Error.InternalServerError(result2.Error, systems.DatabaseError)
	}

	deviceCollection[1] = result2.Value.(*Device)

	device3 := Device{
		GUID:       Helper.GenerateUUID(),
		UUID:       "F50CADC812550C7FC6B3B0A49DBEA480",
		Os:         "Android",
		Model:      "OPPO A1601",
		PushToken:  "cuh2h6UWq-0:APA91bHKfVGhWAoO_x2xAZ2WfdA-Xraw78Eyp58fMsqeJ_jbrt1vkZkTnwKPFvgYLtVP55nXBKseqpo0mqQ7JmK3iPGYFjptQj72HizU4FJdqbQ2AEgvnNxk6gAPeT9eRv9FNQb_TJyB",
		AppVersion: "0.9.9.3-STAGING-BETA",
	}

	result3 := sd.DB.Create(&device3)

	if result3.Error != nil || result3.RowsAffected == 0 {
		return nil, Error.InternalServerError(result3.Error, systems.DatabaseError)
	}

	deviceCollection[2] = result3.Value.(*Device)

	return deviceCollection, nil
}

// DeviceWithUserGUID function used to create sample device including
// user GUID.
func (sd *SampleData) DeviceWithUserGUID(userGUID string) *Device {

	deviceGUID := Helper.GenerateUUID()

	device := Device{
		GUID:       deviceGUID,
		UserGUID:   &userGUID,
		UUID:       "F4EE037F5340586BE625EE7678797DD0",
		Os:         "Android",
		Model:      "Xiaomi MI 3W",
		PushToken:  "cnFBoDLh0W4:APA91bHQLLmjPt24uHFzFNNDg8PgrYodmcNJTqTNv2uWuSpsWnAYd8KHhL0layfGOH3I4mtf46cZzymmw_flehg0CgKNa8vMj-D-vGEPTmL_A71_30pC3OZOanXcc3zK7B7x-9-_A9Lc",
		AppVersion: "0.9.9.3-STAGING-BETA",
	}

	result := sd.DB.Create(&device)

	return result.Value.(*Device)
}

// SmsHistory function used to create sample sms history.
func (sd *SampleData) SmsHistory(event, verificationCode, recipientNo string) (*SmsHistory, *systems.ErrorData) {
	smsHistory := SmsHistory{
		GUID:             Helper.GenerateUUID(),
		Provider:         "moceansms",
		SmsID:            "shoppermate0106234949665476",
		Text:             "Your verification code is " + verificationCode + " - Shoppermate",
		RecipientNo:      recipientNo,
		VerificationCode: verificationCode,
		Event:            event,
	}

	result := sd.DB.Create(&smsHistory)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, Error.InternalServerError(result.Error, systems.DatabaseError)
	}

	return result.Value.(*SmsHistory), nil
}

// User function used to create sample user data for testing database.
func (sd *SampleData) User(phoneNo, email string) *User {
	facebookID := "100013413336774"

	user := User{
		GUID:         Helper.GenerateUUID(),
		Name:         "User 1",
		FacebookID:   &facebookID,
		Email:        email,
		PhoneNo:      phoneNo,
		RegisterBy:   "facebook",
		ReferralCode: "USE24853",
	}

	result := sd.DB.Create(&user)

	return result.Value.(*User)
}

// UserWithCustomWalletAmount function used to create user with custom wallet amount for testing
// database. This useful when you want to test cashout transaction.
func (sd *SampleData) UserWithCustomWalletAmount(phoneNo, email string, walletAmount float64) *User {
	user := User{
		GUID:         Helper.GenerateUUID(),
		Name:         "User 1",
		FacebookID:   nil,
		Email:        email,
		PhoneNo:      phoneNo,
		RegisterBy:   "phone_no",
		ReferralCode: "USE24853",
		Wallet:       walletAmount,
	}

	result := sd.DB.Create(&user)

	return result.Value.(*User)
}

// Users function used to create sample user register using phone number and
// facebook.
func (sd *SampleData) Users() []*User {

	facebookID := "100013413336774"

	user1 := User{
		GUID:         Helper.GenerateUUID(),
		Name:         "User 1",
		FacebookID:   &facebookID,
		Email:        "user1@mediacliq.my",
		PhoneNo:      "60121234567",
		RegisterBy:   "facebook",
		ReferralCode: "USE24853",
	}

	result1 := sd.DB.Create(&user1)

	user2 := User{
		GUID:         Helper.GenerateUUID(),
		Name:         "User 2",
		Email:        "user2@mediacliq.my",
		PhoneNo:      "60171234567",
		RegisterBy:   "phone_no",
		ReferralCode: "USE29563",
	}

	result2 := sd.DB.Create(&user2)

	data := make([]*User, 2)
	data[0] = result1.Value.(*User)
	data[1] = result2.Value.(*User)

	return data
}

// Occasions function used to create sample occasions for test database.
func (sd *SampleData) Occasions() []*Occasion {
	occasion1 := Occasion{
		GUID:   Helper.GenerateUUID(),
		Slug:   "field_trip",
		Name:   "Field Trip",
		Image:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/occasion_images/field_trip.jpg",
		Active: 1,
	}

	result1 := sd.DB.Create(&occasion1)

	occasion2 := Occasion{
		GUID:   Helper.GenerateUUID(),
		Slug:   "travel",
		Name:   "Travel",
		Image:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/occasion_images/travel.jpg",
		Active: 1,
	}

	result2 := sd.DB.Create(&occasion2)

	occasion3 := Occasion{
		GUID:   Helper.GenerateUUID(),
		Slug:   "festive",
		Name:   "Festive",
		Image:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/occasion_images/festive.jpg",
		Active: 0,
	}

	result3 := sd.DB.Create(&occasion3)

	occasion4 := Occasion{
		GUID:   Helper.GenerateUUID(),
		Slug:   "birthday",
		Name:   "Birthday",
		Image:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/occasion_images/birthday.jpg",
		Active: 0,
	}

	result4 := sd.DB.Create(&occasion4)

	data := make([]*Occasion, 4)
	data[0] = result1.Value.(*Occasion)
	data[1] = result2.Value.(*Occasion)
	data[2] = result3.Value.(*Occasion)
	data[3] = result4.Value.(*Occasion)

	return data
}

// ShoppingList function used to create sample shopping list for test database.
func (sd *SampleData) ShoppingList(userGUID, occasionGUID, name string) *ShoppingList {
	shoppingList := ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		OccasionGUID: occasionGUID,
		Name:         name,
	}

	result := sd.DB.Create(&shoppingList)

	return result.Value.(*ShoppingList)
}

// ShoppingLists function used to create sample shopping lists for test database.
func (sd *SampleData) ShoppingLists(userGUID, occasionGUID string) []*ShoppingList {
	shoppingList1 := ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		OccasionGUID: occasionGUID,
		Name:         "Family BBQ",
	}

	result1 := sd.DB.Create(&shoppingList1)

	shoppingList2 := ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		OccasionGUID: occasionGUID,
		Name:         "Birthday Party",
	}

	result2 := sd.DB.Create(&shoppingList2)

	shoppingList3 := ShoppingList{
		GUID:         Helper.GenerateUUID(),
		UserGUID:     userGUID,
		OccasionGUID: occasionGUID,
		Name:         "Wedding Party",
	}

	result3 := sd.DB.Create(&shoppingList3)

	data := make([]*ShoppingList, 3)
	data[0] = result1.Value.(*ShoppingList)
	data[1] = result2.Value.(*ShoppingList)
	data[2] = result3.Value.(*ShoppingList)

	return data
}

// ShoppingListItem function used to create one sample shopping list item for test database.
func (sd *SampleData) ShoppingListItem(userGUID, shoppingListGUID, name string, addedToCart int) *ShoppingListItem {
	shoppingListItem := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             name,
		Quantity:         2,
		Category:         "Others",
		SubCategory:      "Others",
		Remark:           "Buy at Tesco",
		AddedToCart:      addedToCart,
	}

	result := sd.DB.Create(&shoppingListItem)

	return result.Value.(*ShoppingListItem)
}

// ShoppingListItemForDeal function used to create one sample shopping list item when user add
// the deal into cashback for test database.
func (sd *SampleData) ShoppingListItemForDeal(userGUID, shoppingListGUID, name, category, subCategory, dealGUID string,
	cashbackAmount float64, addedToCart int) *ShoppingListItem {

	shoppingListItem := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             name,
		Quantity:         1,
		Category:         category,
		SubCategory:      subCategory,
		AddedFromDeal:    1,
		DealGUID:         &dealGUID,
		CashbackAmount:   &cashbackAmount,
		Remark:           "",
		AddedToCart:      addedToCart,
	}

	result := sd.DB.Create(&shoppingListItem)

	return result.Value.(*ShoppingListItem)
}

// // ShoppingListItems function used to create sample shopping list item for test database.
func (sd *SampleData) ShoppingListItems(userGUID, shoppingListGUID string, addedToCart int) []*ShoppingListItem {
	shoppingListItem1 := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             "Test Shopping List Item 1",
		Quantity:         1,
		Category:         "Others",
		SubCategory:      "Others",
		Remark:           "",
		AddedToCart:      addedToCart,
	}

	result1 := sd.DB.Create(&shoppingListItem1)

	shoppingListItem2 := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             "Test Shopping List Item 2",
		Quantity:         2,
		Category:         "Fresh Food",
		SubCategory:      "Fresh Fruits",
		Remark:           "",
		AddedToCart:      addedToCart,
	}

	result2 := sd.DB.Create(&shoppingListItem2)

	shoppingListItem3 := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             "Test Shopping List Item 3",
		Quantity:         4,
		Category:         "Others",
		SubCategory:      "Others",
		Remark:           "",
		AddedToCart:      addedToCart,
	}

	result3 := sd.DB.Create(&shoppingListItem3)

	shoppingListItem4 := ShoppingListItem{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		ShoppingListGUID: shoppingListGUID,
		Name:             "Test Shopping List Item 4",
		Quantity:         3,
		Category:         "Others",
		SubCategory:      "Others",
		Remark:           "",
		AddedToCart:      addedToCart,
	}

	result4 := sd.DB.Create(&shoppingListItem4)

	data := make([]*ShoppingListItem, 4)
	data[0] = result1.Value.(*ShoppingListItem)
	data[1] = result2.Value.(*ShoppingListItem)
	data[2] = result3.Value.(*ShoppingListItem)
	data[3] = result4.Value.(*ShoppingListItem)

	return data

}

// ShoppingListItemImage function used to create sample shopping list item image for test database.
func (sd *SampleData) ShoppingListItemImage(userGUID, shoppingListGUID, shoppingListItemGUID, imageURL string) *ShoppingListItemImage {
	shoppingListItemImage := ShoppingListItemImage{
		GUID:                 Helper.GenerateUUID(),
		UserGUID:             userGUID,
		ShoppingListGUID:     shoppingListGUID,
		ShoppingListItemGUID: shoppingListItemGUID,
		URL:                  imageURL,
	}

	result := sd.DB.Create(&shoppingListItemImage)

	return result.Value.(*ShoppingListItemImage)
}

// Grocers function used to create sample grocer for test database.
func (sd *SampleData) Grocers() []*Grocer {
	grocer1 := Grocer{
		GUID:   Helper.GenerateUUID(),
		Name:   "Cold Storage Supermarket",
		Email:  "contact@coldstorage.com",
		Status: "publish",
		Img:    "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/grocer_images/42e73d0c-1af0-40e4-871d-887a552da8c0.png",
	}

	result1 := sd.DB.Create(&grocer1)

	grocer2 := Grocer{
		GUID:   Helper.GenerateUUID(),
		Name:   "Jaya Grocer Supermarket",
		Email:  "contact@jayagrocer.com",
		Status: "publish",
		Img:    "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/grocer_images/7bfbe364-d750-4275-9bd9-8faa65e48d3d.jpg",
	}

	result2 := sd.DB.Create(&grocer2)

	grocer3 := Grocer{
		GUID:   Helper.GenerateUUID(),
		Name:   "Ben's Independent Grocer (BIG)",
		Email:  "contact@big.com",
		Status: "publish",
		Img:    "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/grocer_images/77a0de57-54c6-4262-800c-0cb9b0d132bb.png",
	}

	result3 := sd.DB.Create(&grocer3)

	grocer4 := Grocer{
		GUID:   Helper.GenerateUUID(),
		Name:   "Village Grocer",
		Email:  "contact@villagegrocer.com",
		Status: "draft",
		Img:    "https://s3-ap-southeast-1.amazonaws.com/shoppermate-staging/grocer_images/90c01bdd-9e91-4b13-93d6-ab086547f850.jpg",
	}

	result4 := sd.DB.Create(&grocer4)

	grocer5 := Grocer{
		GUID:   Helper.GenerateUUID(),
		Name:   "Tesco",
		Email:  "contact@tesco.com",
		Status: "draft",
		Img:    "https://s3-ap-southeast-1.amazonaws.com/shoppermate-staging/grocer_images/cd24aab1-09a8-4563-9daa-2f88679be818.png",
	}

	result5 := sd.DB.Create(&grocer5)

	data := make([]*Grocer, 5)
	data[0] = result1.Value.(*Grocer)
	data[1] = result2.Value.(*Grocer)
	data[2] = result3.Value.(*Grocer)
	data[3] = result4.Value.(*Grocer)
	data[4] = result5.Value.(*Grocer)

	return data
}

// GrocerLocations function used to create sample grocer location.
// List of Grocer ID
// 1 - Cold Storage
// 2 - Jaya Grocer
// 3 - Ben's Independant Grocer
// 4 - Village Grocer
// 5 - Tesco
func (sd *SampleData) GrocerLocations() []*GrocerLocation {
	grocerLocation1 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 1,
		Name:     "Cold Storage Solaris, Mont Kiara",
		Lat:      3.1762919672722454,
		Lng:      101.65968596935272,
	}

	result1 := sd.DB.Create(&grocerLocation1)

	grocerLocation2 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 1,
		Name:     "Cold Storage Alamanda, Putrajaya ",
		Lat:      2.938787171808547,
		Lng:      101.71214461326599,
	}

	result2 := sd.DB.Create(&grocerLocation2)

	grocerLocation3 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 2,
		Name:     "Jaya Grocer The Intermark",
		Lat:      3.161759217671924,
		Lng:      101.71990424394608,
	}

	result3 := sd.DB.Create(&grocerLocation3)

	grocerLocation4 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 2,
		Name:     "Jaya Grocer Plaza Jelutong",
		Lat:      3.1157962032164037,
		Lng:      101.52900606393814,
	}

	result4 := sd.DB.Create(&grocerLocation4)

	grocerLocation5 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 3,
		Name:     "BIG Publika, Mont Kiara",
		Lat:      3.1710790547708037,
		Lng:      101.66522204875946,
	}

	result5 := sd.DB.Create(&grocerLocation5)

	grocerLocation6 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 3,
		Name:     "BIG Glo Damansara, Damansara",
		Lat:      3.133030561788681,
		Lng:      101.62985175848007,
	}

	result6 := sd.DB.Create(&grocerLocation6)

	grocerLocation7 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 4,
		Name:     "Village Grocer 1MK, Mont Kiara",
		Lat:      3.1659315,
		Lng:      101.65302099999997,
	}

	result7 := sd.DB.Create(&grocerLocation7)

	grocerLocation8 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 4,
		Name:     "Village Grocer Sunway Giza Shopping Mall",
		Lat:      3.1504226477900192,
		Lng:      101.5919092297554,
	}

	result8 := sd.DB.Create(&grocerLocation8)

	grocerLocation9 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 5,
		Name:     "Tesco Kepong Village Mall",
		Lat:      3.193526661982856,
		Lng:      101.63355588912964,
	}

	result9 := sd.DB.Create(&grocerLocation9)

	grocerLocation10 := GrocerLocation{
		GUID:     Helper.GenerateUUID(),
		GrocerID: 5,
		Name:     "Tesco Extra Shah Alam",
		Lat:      3.0717863826292273,
		Lng:      101.53871834278107,
	}

	result10 := sd.DB.Create(&grocerLocation10)

	data := make([]*GrocerLocation, 10)
	data[0] = result1.Value.(*GrocerLocation)
	data[1] = result2.Value.(*GrocerLocation)
	data[2] = result3.Value.(*GrocerLocation)
	data[3] = result4.Value.(*GrocerLocation)
	data[4] = result5.Value.(*GrocerLocation)
	data[5] = result6.Value.(*GrocerLocation)
	data[6] = result7.Value.(*GrocerLocation)
	data[7] = result8.Value.(*GrocerLocation)
	data[8] = result9.Value.(*GrocerLocation)
	data[9] = result10.Value.(*GrocerLocation)

	return data
}

// Categories function used to create sample categories for testing database.
func (sd *SampleData) Categories() []*ItemCategory {
	category1 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/179bd0f6-d3c6-4a6a-a41a-a6485c13eed5.jpg",
		Name: "Baby",
	}

	result1 := sd.DB.Create(&category1)

	category2 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/5a8d71da-6f65-4b38-abf3-c3a3de3e1f61.jpg",
		Name: "Chilled & Frozen",
	}

	result2 := sd.DB.Create(&category2)

	category3 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/a8405c1a-13ff-4053-a8de-be905d915451.jpg",
		Name: "Drinks",
	}

	result3 := sd.DB.Create(&category3)

	category4 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/894b3a22-6e42-485b-a97d-e2e5c04af8a2.jpg",
		Name: "Fresh Food",
	}

	result4 := sd.DB.Create(&category4)

	category5 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/91575eb8-dd95-4530-84fc-276619e5d194.jpg",
		Name: "Grocery",
	}

	result5 := sd.DB.Create(&category5)

	category6 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/e3191cde-4f93-4f17-9961-67b233749621.jpg",
		Name: "Health & Beauty",
	}

	result6 := sd.DB.Create(&category6)

	category7 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/8fe2fd22-1afe-45ee-a968-4f777354f2d1.jpg",
		Name: "Household",
	}

	result7 := sd.DB.Create(&category7)

	category8 := ItemCategory{
		GUID: Helper.GenerateUUID(),
		Img:  "https://s3-ap-southeast-1.amazonaws.com/shoppermate/category_images/b21792b6-b316-4199-a0bc-7dc506f7a57d.jpg",
		Name: "Non-Halal",
	}

	result8 := sd.DB.Create(&category8)

	data := make([]*ItemCategory, 8)
	data[0] = result1.Value.(*ItemCategory)
	data[1] = result2.Value.(*ItemCategory)
	data[2] = result3.Value.(*ItemCategory)
	data[3] = result4.Value.(*ItemCategory)
	data[4] = result5.Value.(*ItemCategory)
	data[5] = result6.Value.(*ItemCategory)
	data[6] = result7.Value.(*ItemCategory)
	data[7] = result8.Value.(*ItemCategory)

	return data
}

// Subcategories function used to create sample subcategories for testing database.
func (sd *SampleData) Subcategories() []*ItemSubCategory {
	subcategory1 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Carbonated Drink",
	}

	result1 := sd.DB.Create(&subcategory1)

	subcategory2 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Chocolate & Sweets",
	}

	result2 := sd.DB.Create(&subcategory2)

	subcategory3 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Milk",
	}

	result3 := sd.DB.Create(&subcategory3)

	subcategory4 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Flavoured Beverages",
	}

	result4 := sd.DB.Create(&subcategory4)

	subcategory5 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Coffee",
	}

	result5 := sd.DB.Create(&subcategory5)

	subcategory6 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Hair Care",
	}

	result6 := sd.DB.Create(&subcategory6)

	subcategory7 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Butter & Margarine",
	}

	result7 := sd.DB.Create(&subcategory7)

	subcategory8 := ItemSubCategory{
		GUID: Helper.GenerateUUID(),
		Name: "Jam, Speads & Honey",
	}

	result8 := sd.DB.Create(&subcategory8)

	data := make([]*ItemSubCategory, 8)
	data[0] = result1.Value.(*ItemSubCategory)
	data[1] = result2.Value.(*ItemSubCategory)
	data[2] = result3.Value.(*ItemSubCategory)
	data[3] = result4.Value.(*ItemSubCategory)
	data[4] = result5.Value.(*ItemSubCategory)
	data[5] = result6.Value.(*ItemSubCategory)
	data[6] = result7.Value.(*ItemSubCategory)
	data[7] = result8.Value.(*ItemSubCategory)

	return data
}

// Generics function used to create sample generics data for testing database.
// List of Category ID               List Of Subcategory ID
// 1 - Baby                          1 - Carbonated Drinks
// 2 - Chilled & Frozen              2 - Chocolate & Sweets
// 3 - Drinks                        3 - Milk
// 4 - Fresh Food                    4 - Flavoured Beverages
// 5 - Grocery                       5 - Coffee
// 6 - Health & Beauty               6 - Hair Care
// 7 - Household                     7 - Butter & Margarine
// 8 - Non-Halal                     8 - Jam, Speads & Honey
func (sd *SampleData) Generics() []*Generic {
	generic1 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    3,
		SubcategoryID: 1,
		Name:          "Cola",
	}

	result1 := sd.DB.Create(&generic1)

	generic2 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    5,
		SubcategoryID: 2,
		Name:          "Chocolates",
	}

	result2 := sd.DB.Create(&generic2)

	generic3 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    2,
		SubcategoryID: 3,
		Name:          "Chocolate Milk",
	}

	result3 := sd.DB.Create(&generic3)

	generic4 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    3,
		SubcategoryID: 4,
		Name:          "Chocolate Drink",
	}

	result4 := sd.DB.Create(&generic4)

	generic5 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    3,
		SubcategoryID: 5,
		Name:          "Instant Coffee",
	}

	result5 := sd.DB.Create(&generic5)

	generic6 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    6,
		SubcategoryID: 6,
		Name:          "Hair Shampoo",
	}

	result6 := sd.DB.Create(&generic6)

	generic7 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    6,
		SubcategoryID: 6,
		Name:          "Hair Conditioner",
	}

	result7 := sd.DB.Create(&generic7)

	generic8 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    2,
		SubcategoryID: 7,
		Name:          "Margerine",
	}

	result8 := sd.DB.Create(&generic8)

	generic9 := Generic{
		GUID:          Helper.GenerateUUID(),
		CategoryID:    5,
		SubcategoryID: 8,
		Name:          "Margarine",
	}

	result9 := sd.DB.Create(&generic9)

	data := make([]*Generic, 9)
	data[0] = result1.Value.(*Generic)
	data[1] = result2.Value.(*Generic)
	data[2] = result3.Value.(*Generic)
	data[3] = result4.Value.(*Generic)
	data[4] = result5.Value.(*Generic)
	data[5] = result6.Value.(*Generic)
	data[6] = result7.Value.(*Generic)
	data[7] = result8.Value.(*Generic)
	data[8] = result9.Value.(*Generic)

	return data
}

// Items function used to create sample items for testing database.
// List of Category ID               List Of Subcategory ID                List Of Generic ID
// 1 - Baby                          1 - Carbonated Drinks                 1 - Cola
// 2 - Chilled & Frozen              2 - Chocolate & Sweets                2 - Chocolates
// 3 - Drinks                        3 - Milk                              3 - Chocolate Milk
// 4 - Fresh Food                    4 - Flavoured Beverages               4 - Chocolate Drink
// 5 - Grocery                       5 - Coffee                            5 - Instant Coffee
// 6 - Health & Beauty               6 - Hair Care                         6 - Hair Shampoo
// 7 - Household                     7 - Butter & Margarine                7 - Hair Conditioner
// 8 - Non-Halal                     8 - Jam, Speads & Honey               8 - Margerine
//                                                                         9 - Margarine
func (sd *SampleData) Items() []*Item {
	genericID1 := 1
	genericID2 := 2
	genericID3 := 3
	genericID4 := 4
	genericID5 := 5
	genericID6 := 6
	genericID7 := 7
	genericID8 := 8

	item1 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID1,
		Name:          "Coca-cola Vanilla 1.5L",
		CategoryID:    3,
		SubcategoryID: 1,
		Remarks:       "",
	}

	result1 := sd.DB.Create(&item1)

	item2 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID1,
		Name:          "Coca-cola 500mL",
		CategoryID:    3,
		SubcategoryID: 1,
		Remarks:       "",
	}

	result2 := sd.DB.Create(&item2)

	item3 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID2,
		Name:          "Cadbury Dairy Milk Chocolate 100pcs",
		CategoryID:    5,
		SubcategoryID: 2,
		Remarks:       "",
	}

	result3 := sd.DB.Create(&item3)

	item4 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID3,
		Name:          "Dutch Lady Chocolate Milk 1L",
		CategoryID:    3,
		SubcategoryID: 3,
		Remarks:       "",
	}

	result4 := sd.DB.Create(&item4)

	item5 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID2,
		Name:          "Kit Kat Chunky Peanut Butter",
		CategoryID:    5,
		SubcategoryID: 2,
		Remarks:       "",
	}

	result5 := sd.DB.Create(&item5)

	item6 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID4,
		Name:          "Milo Activ-go Tin 1.5kg",
		CategoryID:    3,
		SubcategoryID: 4,
		Remarks:       "",
	}

	result6 := sd.DB.Create(&item6)

	item7 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID5,
		Name:          "Nescafe Gold Rich Aroma Pure Soluble Coffee 200g",
		CategoryID:    3,
		SubcategoryID: 5,
		Remarks:       "",
	}

	result7 := sd.DB.Create(&item7)

	item8 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID6,
		Name:          "Pantene Pro-v Aqua Pure Shampoo 750mL",
		CategoryID:    6,
		SubcategoryID: 6,
		Remarks:       "",
	}

	result8 := sd.DB.Create(&item8)

	item9 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID7,
		Name:          "Pantene Pro-v Daily Moisture Repair Conditioner 670mL",
		CategoryID:    6,
		SubcategoryID: 6,
		Remarks:       "",
	}

	result9 := sd.DB.Create(&item9)

	item10 := Item{
		GUID:          Helper.GenerateUUID(),
		GenericID:     &genericID8,
		Name:          "Naturel Lite Spread 500g",
		CategoryID:    2,
		SubcategoryID: 7,
		Remarks:       "",
	}

	result10 := sd.DB.Create(&item10)

	data := make([]*Item, 10)
	data[0] = result1.Value.(*Item)
	data[1] = result2.Value.(*Item)
	data[2] = result3.Value.(*Item)
	data[3] = result4.Value.(*Item)
	data[4] = result5.Value.(*Item)
	data[5] = result6.Value.(*Item)
	data[6] = result7.Value.(*Item)
	data[7] = result8.Value.(*Item)
	data[8] = result9.Value.(*Item)
	data[9] = result10.Value.(*Item)

	return data
}

// Deals function used to create sample deal data in testing database.
// Grocer Exclusive ID
// 1 -
// 2 -
// 3 -
// 4 -
// 5 -
func (sd *SampleData) Deals() []*Ads {
	dealStartDate := time.Now().UTC().Add(time.Hour * -8)
	dealEndDate := time.Now().UTC().Add(time.Hour * 24)

	grocerExclusive1 := 1
	grocerExclusive2 := 2
	grocerExclusive3 := 3
	grocerExclusive4 := 4
	grocerExclusive5 := 5

	deal1 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/cffb3ae5-8530-400b-97b5-8b6af2a36691.jpg",
		FrontName:       "Coca-Cola Vanilla 1.5L, Get it now!",
		Name:            "Coca-Cola Vanilla 1.5L",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          1,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        2,
		CashbackAmount:  1.00,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive1,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result1 := sd.DB.Create(&deal1)

	deal2 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/8d9a2ae0-5f58-4a2d-8128-526b2ea886ae.jpg",
		FrontName:       "Coca-Cola 500ml - Any Variant",
		Name:            "Coca-Cola 500ml",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          2,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        1,
		CashbackAmount:  0.50,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: nil,
		Terms:           "",
	}

	result2 := sd.DB.Create(&deal2)

	deal3 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/77b20384-bd7d-4717-959e-b7e2d3aeeea5.jpg",
		FrontName:       "Cadbury Dairy Milk Chocolate",
		Name:            "Cadbury Dairy Milk Chocolate - 100pcs",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          3,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        2,
		CashbackAmount:  1.20,
		Quota:           2,
		Status:          "publish",
		GrocerExclusive: nil,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result3 := sd.DB.Create(&deal3)

	deal4 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/00c6f2f8-1e6b-40ca-89c7-f3c4d50f6916.jpg",
		FrontName:       "Dutch Lady Chocolate Milk 1L",
		Name:            "Dutch Lady Chocolate Milk 1L",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          4,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  0.75,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive5,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result4 := sd.DB.Create(&deal4)

	deal5 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/0ea9db7a-a056-4174-942f-4837d5824874.jpg",
		FrontName:       "Kit Kat Chunky - Any Variant",
		Name:            "Kit Kat Chunky",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          5,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  1.00,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: nil,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result5 := sd.DB.Create(&deal5)

	deal6 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/3f4f2093-2c8a-41ca-8bee-d5a2c6de55bc.jpg",
		FrontName:       "Milo Activ-Go Tin 1.5Kg",
		Name:            "Milo Activ-Go Tin 1.5Kg",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          6,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  1.50,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive2,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result6 := sd.DB.Create(&deal6)

	deal7 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate/deal_images/9126c1fc-d455-46b0-9f67-a03d86e51639.jpg",
		FrontName:       "Nescafe Gold Rich Aroma Pure Soluble Coffee 200g",
		Name:            "Nescafe Gold Rich Aroma Pure Soluble Coffee 200g",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          7,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  2.00,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive3,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result7 := sd.DB.Create(&deal7)

	deal8 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-staging/deal_images/b0cdb4ec-0fe9-4e32-b55c-796718cf0993.jpg",
		FrontName:       "Pantene Pro-V Aqua Pure Shampoo 750ml",
		Name:            "Pantene Pro-V Aqua Pure Shampoo 750ml",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          8,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  1.45,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: nil,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result8 := sd.DB.Create(&deal8)

	deal9 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/52c278ac-5d59-46e3-bd64-1a9daf44135f.jpg",
		FrontName:       "Pantene Pro-V  Conditioner 670ml - Any Variant",
		Name:            "Pantene Pro-V  Conditioner 670ml - Any Variant",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          9,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  2.30,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive4,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result9 := sd.DB.Create(&deal9)

	deal10 := Ads{
		GUID:            Helper.GenerateUUID(),
		Img:             "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_images/bd65afdc-0e10-4781-8407-d14adafe6b13.jpg",
		FrontName:       "Naturel Lite Spread 500g",
		Name:            "Naturel Lite Spread 500g",
		Body:            "Cashback valid for purchase of 1 unit of the item from any of the participating grocer",
		AdvertiserID:    1,
		CampaignID:      1,
		ItemID:          10,
		PositiveTag:     "",
		NegativeTag:     "",
		Type:            "deal",
		StartDate:       dealStartDate,
		EndDate:         dealEndDate,
		Time:            "",
		RefreshPeriod:   1,
		Perlimit:        3,
		CashbackAmount:  1.00,
		Quota:           20,
		Status:          "publish",
		GrocerExclusive: &grocerExclusive4,
		Terms: `Deal must first be ""Added to list"" before purchase is made.
				Purchase is made on the right item, varient, quantity and size.
				Purchase (Date of receipt) is made within the stipulated deal period.
				Purchase is made in one of the listed participating grocer and outlet.
				Cashback redemption is made within 7 days from the date of the receipt via the ShopperMate App.
				These terms are supplementary to the ShopperMate Term of Use found here: http://www.shoppermate.com/tou`,
	}

	result10 := sd.DB.Create(&deal10)

	data := make([]*Ads, 10)
	data[0] = result1.Value.(*Ads)
	data[1] = result2.Value.(*Ads)
	data[2] = result3.Value.(*Ads)
	data[3] = result4.Value.(*Ads)
	data[4] = result5.Value.(*Ads)
	data[5] = result6.Value.(*Ads)
	data[6] = result7.Value.(*Ads)
	data[7] = result8.Value.(*Ads)
	data[8] = result9.Value.(*Ads)
	data[9] = result10.Value.(*Ads)

	return data
}

// DealCashback function used to create sample Deal Cashback for test database.
func (sd *SampleData) DealCashback(userGUID, shoppingListGUID, dealGUID string, dealCashbackTransactionGUID *string) *DealCashback {
	dealCashback := &DealCashback{
		GUID:                        Helper.GenerateUUID(),
		UserGUID:                    userGUID,
		ShoppingListGUID:            shoppingListGUID,
		DealGUID:                    dealGUID,
		DealCashbackTransactionGUID: dealCashbackTransactionGUID,
	}

	result := sd.DB.Create(dealCashback)

	return result.Value.(*DealCashback)
}

// DealCashbackTransactionWithPendingCleaningStatus function used to create sample Deal Cashback Transaction
// with pending status for test database. Pending status means the verification date, remark title and
// remark body is nil and the status must be 'pending cleaning'
func (sd *SampleData) DealCashbackTransactionWithPendingCleaningStatus(dealCashbackGUID, userGUID, transactionGUID string) *DealCashbackTransaction {
	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		TransactionGUID:  transactionGUID,
		ReceiptURL:       "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_cashback_receipts/test_receipt.jpg",
		VerificationDate: nil,
		RemarkTitle:      nil,
		RemarkBody:       nil,
		Status:           "pending cleaning",
	}

	result := sd.DB.Create(dealCashbackTransaction)

	return result.Value.(*DealCashbackTransaction)
}

// DealCashbackTransactionWithPendingApprovalStatus function used to create sample Deal Cashback Transaction
// with completed status for test database. Completed status means the verification date can't be empty
// and status value must be 'pending approval'
func (sd *SampleData) DealCashbackTransactionWithPendingApprovalStatus(dealCashbackGUID, userGUID, transactionGUID string) *DealCashbackTransaction {
	verificationDate := time.Now().UTC().Add(time.Hour * -8)

	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		TransactionGUID:  transactionGUID,
		ReceiptURL:       "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_cashback_receipts/test_receipt.jpg",
		VerificationDate: &verificationDate,
		RemarkTitle:      nil,
		RemarkBody:       nil,
		Status:           "pending approval",
	}

	result := sd.DB.Create(dealCashbackTransaction)

	return result.Value.(*DealCashbackTransaction)
}

// DealCashbackTransactionWithCompletedStatus function used to create sample Deal Cashback Transaction
// with completed status for test database. Completed status means the verification date can't be empty
// and status value must be 'completed'
func (sd *SampleData) DealCashbackTransactionWithCompletedStatus(dealCashbackGUID, userGUID, transactionGUID string) *DealCashbackTransaction {
	verificationDate := time.Now().UTC().Add(time.Hour * -8)

	dealCashbackTransaction := &DealCashbackTransaction{
		GUID:             Helper.GenerateUUID(),
		UserGUID:         userGUID,
		TransactionGUID:  transactionGUID,
		ReceiptURL:       "https://s3-ap-southeast-1.amazonaws.com/shoppermate-test/deal_cashback_receipts/test_receipt.jpg",
		VerificationDate: &verificationDate,
		RemarkTitle:      nil,
		RemarkBody:       nil,
		Status:           "completed",
	}

	result := sd.DB.Create(dealCashbackTransaction)

	return result.Value.(*DealCashbackTransaction)
}

// TransactionStatuses function used to create sample transaction statuses for test database.
func (sd *SampleData) TransactionStatuses() []*TransactionStatus {
	pendingStatus := &TransactionStatus{
		GUID: Helper.GenerateUUID(),
		Slug: "pending",
		Name: "pending",
	}

	result1 := sd.DB.Create(pendingStatus)

	partialSuccessStatus := &TransactionStatus{
		GUID: Helper.GenerateUUID(),
		Slug: "partial_success",
		Name: "partial success",
	}

	result2 := sd.DB.Create(partialSuccessStatus)

	approvedStatus := &TransactionStatus{
		GUID: Helper.GenerateUUID(),
		Slug: "approved",
		Name: "approved",
	}

	result3 := sd.DB.Create(approvedStatus)

	rejectStatus := &TransactionStatus{
		GUID: Helper.GenerateUUID(),
		Slug: "reject",
		Name: "reject",
	}

	result4 := sd.DB.Create(rejectStatus)

	data := make([]*TransactionStatus, 4)
	data[0] = result1.Value.(*TransactionStatus)
	data[1] = result2.Value.(*TransactionStatus)
	data[2] = result3.Value.(*TransactionStatus)
	data[3] = result4.Value.(*TransactionStatus)

	return data
}

// TransactionTypes function used to create sample transaction types for test database.
func (sd *SampleData) TransactionTypes() []*TransactionType {
	referralCashbackType := &TransactionType{
		GUID: Helper.GenerateUUID(),
		Slug: "referral_cashback",
		Name: "Referral Cashback",
	}

	result1 := sd.DB.Create(referralCashbackType)

	dealRedemptionType := &TransactionType{
		GUID: Helper.GenerateUUID(),
		Slug: "deal_redemption",
		Name: "Deal Redemption",
	}

	result2 := sd.DB.Create(dealRedemptionType)

	cashoutType := &TransactionType{
		GUID: Helper.GenerateUUID(),
		Slug: "cashout",
		Name: "Cashout",
	}

	result3 := sd.DB.Create(cashoutType)

	data := make([]*TransactionType, 3)
	data[0] = result1.Value.(*TransactionType)
	data[1] = result2.Value.(*TransactionType)
	data[2] = result3.Value.(*TransactionType)

	return data
}

// Transaction function used to create sample transaction for testing database.
func (sd *SampleData) Transaction(userGUID, transactionTypeGUID, transactionStatusGUID string, readStatus int, totalAmount float64) *Transaction {
	transaction := &Transaction{
		GUID:                  Helper.GenerateUUID(),
		UserGUID:              userGUID,
		TransactionTypeGUID:   transactionTypeGUID,
		TransactionStatusGUID: transactionStatusGUID,
		ReadStatus:            0,
		ReferenceID:           Helper.GenerateUniqueShortID(),
		TotalAmount:           totalAmount,
	}

	result := sd.DB.Create(transaction)

	return result.Value.(*Transaction)
}
