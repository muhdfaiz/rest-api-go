package v1_1

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/jinzhu/gorm"
)

// SampleData will create all data needed during testing.
type SampleData struct {
	DB *gorm.DB
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
func (sd *SampleData) DeviceWithUserGUID(userGUID string) interface{} {

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

	return result.Value
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
		RegisterBy:   "phone_no",
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
