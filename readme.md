## Sample API in Go. Took from my abandoned project.

### Development Enviroment

### Prerequisite

- Go Languange 1.7 above, MariaDB, Glide (Package Management For Go, GIT, supervisor

#### Install MariaDB 10.x

```
Refer here - https://www.linuxbabe.com/mariadb/install-mariadb-10-1-ubuntu14-04-15-10
```

#### Install Git

```
sudo apt-get update
sudo apt-get install git
```

#### Install Go Language 1.7.x

- Download Go Language 1.7.x

```
sudo apt-get update
sudo wget https://storage.googleapis.com/golang/go1.7.4.linux-amd64.tar.gz
```

- Extract Go Language 1.7.x

```
sudo tar -xvf go1.7.4.linux-amd64.tar.gz
sudo mv go /usr/local
```

#### Setup Go Environment.

- In Linux, edit file `~/.profile` or `~/bash_profile` and include 3 environment variables below. If you are using zsh shell you can edit file `.zshrc` If you are using Mac edit this file `~/.bash_profile`

- Set GOROOT (location when Go package is installed on your system

```
export GOROOT=/usr/local/go
```

- Set GOPATH. Location of your project path. For example 

```
export GOPATH=$HOME/golang
```

- Set PATH variable to access go binary system wide.

```
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

#### Verify Installation

- Check Go Version

```
go version
```

- Verify all environment variable. Make sure GOROOT and GOPATH not empty and set to the correct folder.

```
go env
```

#### Install Glide - Package Management for Go (https://github.com/Masterminds/glide

- Install Glide

For Ubuntu

```
sudo add-apt-repository ppa:masterminds/glide && sudo apt-get update
sudo apt-get install glide
```

For Mac Os X
```
brew install glide
```

- Clone Repository

- Go to project path. Install package dependencies using Glide. Glide will install all package dependencies into this folder `~/golang/src/vendor`.

- Refer file `glide.yaml` and `glide.lock` located in root project path to see what package glide will install.

```
glide install
```

- Create new .env file and copy the content from .env.example file in root directory. Update all setting in .env file.

- Go to project root directory and run the file.
```
go run api.go 
```

### Project Documentation

#### Library Used
- ORM: [GORM](http://jinzhu.me/gorm/)
- Framework: [GIN](https://github.com/gin-gonic/gin)
- Validator: [Go Playground Validator v8](https://github.com/go-playground/validator)

#### Routes

- File Location

```
application/v1/routes.go 
application/v11/routes.go 
```

- List Of Endpoint And Associated Handler

##### API Version V1

```
POST   /v1/devices                                                                             --> /application/v1/device_handler.go (DeviceHandler.Create)
PATCH  /v1/devices/:uuid                                                                       --> /application/v1/device_handler.go (DeviceHandler.Update)
POST   /v1/users                                                                               --> /application/v1/user_handler.go (UserHandler.Create)
POST   /v1/sms                                                                                 --> /application/v1/sms_handler.go (SmsHandler.Send)
POST   /v1/sms/verifications                                                                   --> /application/v1/sms_handler.go (SmsHandler.Verify)
POST   /v1/auth/login/phone                                                                    --> /application/v1/auth_handler.go (AuthHandler.LoginViaPhone)
POST   /v1/auth/login/facebook                                                                 --> /application/v1/auth_handler.go (AuthHandler.LoginViaFacebook)
GET    /v1/shopping_lists/occasions                                                            --> /application/v1/occasion_handler.go (OccasionHandler.Index)
GET    /v1/shopping_lists/items                                                                --> /application/v1/item_handler.go (ItemHandler.Index)
GET    /v1/shopping_lists/items/categories                                                     --> /application/v1/item_category_handler.go (ItemCategoryHandler.ViewAll)
GET    /v1/generics                                                                            --> /application/v1/generic_handler.go (GenericHandler.ViewAll)
GET    /v1/deals                                                                               --> /application/v1/deal_handler.go (DealHandler.ViewAllForGuestUser)
GET    /v1/shopping_list_samples                                                               --> /application/v1/default_shopping_list_handler.go DefaultShoppingListHandler.ViewAll)
GET    /v1/settings                                                                            --> /application/v1/setting_handler.go (SettingHandler.ViewAll)
PATCH  /v1/users/:guid                                                                         --> /application/v1/user_handler.go (UserHandler.Update)
GET    /v1/users/:guid                                                                         --> /application/v1/user_handler.go (UserHandler.View)
DELETE /v1/devices/:uuid                                                                       --> /application/v1/device_handler.go (DeviceHandler.Delete)
GET    /v1/auth/refresh                                                                        --> /application/v1/auth_handler.go (AuthHandler.Refresh)
GET    /v1/auth/logout                                                                         --> /application/v1/auth_handler.go (AuthHandler.Logout)
GET    /v1/users/:guid/shopping_lists                                                          --> /application/v1/shopping_list_handler.go (ShoppingListHandler.View)
POST   /v1/users/:guid/shopping_lists                                                          --> /application/v1/shopping_list_handler.go (ShoppingListHandler.Create)
PATCH  /v1/users/:guid/shopping_lists/:shopping_list_guid                                      --> /application/v1/shopping_list_handler.go (ShoppingListHandler.Update)
DELETE /v1/users/:guid/shopping_lists/:shopping_list_guid                                      --> /application/v1/shopping_list_handler.go (ShoppingListHandler.Delete)
GET    /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                     --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.View)
GET    /v1/users/:guid/shopping_lists/:shopping_list_guid/items                                --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.ViewAll)
POST   /v1/users/:guid/shopping_lists/:shopping_list_guid/items                                --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.Create)
PATCH  /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                     --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.Update)
PATCH  /v1/users/:guid/shopping_lists/:shopping_list_guid/items                                --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.UpdateAll)
DELETE /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                     --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.Delete)
DELETE /v1/users/:guid/shopping_lists/:shopping_list_guid/items                                --> /application/v1/shopping_list_item_handler.go (ShoppingListItemHandler.DeleteAll)
GET    /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guid  --> /application/v1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.View)
POST   /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images              --> /application/v1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.Create)
DELETE /v1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guids --> /application/v1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.Delete)
GET    /v1/users/:guid/deals                                                                   --> /application/v1/deal_handler.go (DealHandler.ViewAllForRegisteredUser)
GET    /v1/deals/:deal_guid                                                                    --> /application/v1/deal_handler.go (DealHandler.View)
GET    /v1/users/:guid/deals/categories                                                        --> /application/v1/deal_handler.go (DealHandler.ViewAndGroupByCategory)
GET    /v1/users/:guid/deals/categories/:category_guid                                         --> /application/v1/deal_handler.go (DealHandler.ViewByCategory)
GET    /v1/users/:guid/deals/categories/:category_guid/subcategories                           --> /application/v1/deal_handler.go (DealHandler.ViewByCategoryAndGroupBySubCategory)
GET    /v1/users/:guid/deals/subcategories/:subcategory_guid                                   --> /application/v1/deal_handler.go (DealHandler.ViewBySubCategory)
GET    /v1/users/:guid/deals/grocers                                                           --> /application/v1/grocer_handler.go (GrocerHandler.GetAllGrocersThatContainDeals)
GET    /v1/users/:guid/deals/grocers/:grocer_guid/categories                                   --> /application/v1/item_category_handler.go (ItemCategoryHandlerViewGrocerCategoriesThoseHaveDealsIncludingDeals)
GET    /v1/users/:guid/deals/grocers/:grocer_guid/categories/:category_guid                    --> /application/v1/deal_handler.go (DealHandler.ViewByGrocerAndCategory)
GET    /v1/users/:guid/featured_deals                                                          --> /application/v1/event_handler.go (EventHandler.ViewAll)
POST   /v1/users/:guid/deal_cashbacks                                                          --> /application/v1/deal_cashback_handler.go (DealCashbackHandler.Create)
GET    /v1/users/:guid/deal_cashbacks/shopping_lists/:shopping_list_guid                       --> /application/v1/deal_cashback_handler.go (DealCashbackHandler.ViewByShoppingList)
GET    /v1/users/:guid/deal_cashbacks/deals/:deal_guid                                         --> /application/v1/deal_cashback_handler.go (DealCashbackHandler.ViewByUserAndDealGroupByShoppingList)
POST   /v1/users/:guid/transactions/deal_cashback_transactions                                 --> /application/v1/deal_cashback_transaction_handler.go (DealCashbackTransactionHandler.Create)
GET    /v1/users/:guid/transactions                                                            --> /application/v1/transaction_handler.go (TransactionHandler.ViewUserTransactions)
GET    /v1/users/:guid/transactions/:transaction_guid/deal_cashback_transactions               --> /application/v1/transaction_handler.go  (TransactionHandler.ViewDealCashbackTransaction)
GET    /v1/users/:guid/transactions/:transaction_guid/cashout_transactions                     --> /application/v1/transaction_handler.go  (TransactionHandler.ViewCashoutTransaction)
GET    /v1/users/:guid/transactions/:transaction_guid/referral_cashback_transactions           --> /application/v1/transaction_handler.go (TransactionHandler.ViewReferralCashbackTransaction)
POST   /v1/users/:guid/transactions/cashout_transactions                                       --> /application/v1/cashout_transaction_handler.go (CashoutTransactionHandler.Create)
```

##### API Version V1.1

```
POST   /v1_1/devices                                                                           --> /application/v1_1/device_handler.go (DeviceHandler.Create)
PATCH  /v1_1/devices/:uuid                                                                     --> /application/v1_1/device_handler.go (DeviceHandler.Update)
POST   /v1_1/users                                                                             --> /application/v1_1/user_handler.go (UserHandler.Create)
POST   /v1_1/sms                                                                               --> /application/v1_1/sms_handler.go (SmsHandler.Send)
POST   /v1_1/sms/verifications                                                                 --> /application/v1_1/sms_handler.go (SmsHandler.Verify)
POST   /v1_1/auth/login/phone                                                                  --> /application/v1_1/auth_handler.go (AuthHandler.LoginViaPhone)
POST   /v1_1/auth/login/facebook                                                               --> /application/v1_1/auth_handler.go (AuthHandler.LoginViaFacebook)
GET    /v1_1/shopping_lists/occasions                                                          --> /application/v1_1/occasion_handler.go (OccasionHandler.Index)
GET    /v1_1/shopping_lists/items                                                              --> /application/v1_1/item_handler.go (ItemHandler.Index)
GET    /v1_1/shopping_lists/items/categories                                                   --> /application/v1_1/item_category_handler.go (ItemCategoryHandler.ViewAll)
GET    /v1_1/generics                                                                          --> /application/v1_1/generic_handler.go (GenericHandler.ViewAll)
GET    /v1_1/deals                                                                             --> /application/v1_1/deal_handler.go (DealHandler.ViewAllForGuestUser)
GET    /v1_1/shopping_list_samples                                                             --> /application/v1_1/default_shopping_list_handler.go DefaultShoppingListHandler.ViewAll)
GET    /v1_1/settings                                                                          --> /application/v1_1/setting_handler.go (SettingHandler.ViewAll)
GET    /v1_1/device/:device_uuid/notifications                                                 --> /application/v1_1/notification_handler.go (NotificationHandler.ViewNotificationForGuest)
PATCH  /v1_1/users/:guid                                                                       --> /application/v1_1/user_handler.go (UserHandler.Update)
GET    /v1_1/users/:guid                                                                       --> /application/v1_1/user_handler.go (UserHandler.View)
DELETE /v1_1/devices/:uuid                                                                     --> /application/v1_1/device_handler.go (DeviceHandler.Delete)
GET    /v1_1/auth/refresh                                                                      --> /application/v1_1/auth_handler.go (AuthHandler.Refresh)
GET    /v1_1/auth/logout                                                                       --> /application/v1_1/auth_handler.go (AuthHandler.Logout)
GET    /v1_1/users/:guid/shopping_lists                                                        --> /application/v1_1/shopping_list_handler.go (ShoppingListHandler.View)
POST   /v1_1/users/:guid/shopping_lists                                                        --> /application/v1_1/shopping_list_handler.go (ShoppingListHandler.Create)
PATCH  /v1_1/users/:guid/shopping_lists/:shopping_list_guid                                    --> /application/v1_1/shopping_list_handler.go (ShoppingListHandler.Update)
DELETE /v1_1/users/:guid/shopping_lists/:shopping_list_guid                                    --> /application/v1_1/shopping_list_handler.go (ShoppingListHandler.Delete)
GET    /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                   --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.View)
GET    /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items                              --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.ViewAll)
POST   /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items                              --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.Create)
PATCH  /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                   --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.Update)
PATCH  /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items                              --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.UpdateAll)
DELETE /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid                   --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.Delete)
DELETE /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items                              --> /application/v1_1/shopping_list_item_handler.go (ShoppingListItemHandler.DeleteAll)
GET    /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_gui --> /application/v1_1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.View)
POST   /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images            --> /application/v1_1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.Create)
DELETE /v1_1/users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_gui --> /application/v1_1/shopping_list_item_image_handler.go ShoppingListItemImageHandler.Delete)
GET    /v1_1/users/:guid/deals                                                                 --> /application/v1_1/deal_handler.go (DealHandler.ViewAllForRegisteredUser)
GET    /v1_1/deals/:deal_guid                                                                  --> /application/v1_1/deal_handler.go (DealHandler.View)
GET    /v1_1/users/:guid/deals/categories                                                      --> /application/v1_1/deal_handler.go (DealHandler.ViewAndGroupByCategory)
GET    /v1_1/users/:guid/deals/categories/:category_guid                                       --> /application/v1_1/deal_handler.go (DealHandler.ViewByCategory)
GET    /v1_1/users/:guid/deals/categories/:category_guid/subcategories                         --> /application/v1_1/deal_handler.go (DealHandler.ViewByCategoryAndGroupBySubCategory)
GET    /v1_1/users/:guid/deals/subcategories/:subcategory_guid                                 --> /application/v1_1/deal_handler.go (DealHandler.ViewBySubCategory)
GET    /v1_1/users/:guid/deals/grocers                                                         --> /application/v1_1/grocer_handler.go (GrocerHandler.GetAllGrocersThatContainDeals)
GET    /v1_1/users/:guid/deals/grocers/:grocer_guid/categories                                 --> /application/v1_1/item_category_handler.go (ItemCategoryHandler.ViewGrocerCategoriesThoseHaveDealsIncludingDeals)
GET    /v1_1/users/:guid/deals/grocers/:grocer_guid/categories/:category_guid                  --> /application/v1_1/deal_handler.go (DealHandler.ViewByGrocerAndCategory)
GET    /v1_1/users/:guid/featured_deals                                                        --> /application/v1_1/event_handler.go (EventHandler.ViewAll)
POST   /v1_1/users/:guid/deal_cashbacks                                                        --> /application/v1_1/deal_cashback_handler.go (DealCashbackHandler.Create)
GET    /v1_1/users/:guid/deal_cashbacks/shopping_lists/:shopping_list_guid                     --> /application/v1_1/deal_cashback_handler.go (DealCashbackHandler.ViewByShoppingList)
GET    /v1_1/users/:guid/deal_cashbacks                                                        --> /application/v1_1/deal_cashback_handler.go (DealCashbackHandler.ViewByUserAndGroupByShoppingList)
GET    /v1_1/users/:guid/deal_cashbacks/deals/:deal_guid                                       --> /application/v1_1/deal_cashback_handler.go (DealCashbackHandler.ViewByUserAndDealGroupByShoppingList)
POST   /v1_1/users/:guid/transactions/deal_cashback_transactions                               --> /application/v1_1/deal_cashback_transaction_handler.go (DealCashbackTransactionHandler.Create)
GET    /v1_1/users/:guid/transactions                                                          --> /application/v1_1/transaction_handler.go (TransactionHandler.ViewUserTransactions)
GET    /v1_1/users/:guid/transactions/:transaction_guid/deal_cashback_transactions             --> /application/v1_1/transaction_handler.go (TransactionHandler.ViewDealCashbackTransaction)
GET    /v1_1/users/:guid/transactions/:transaction_guid/cashout_transactions                   --> /application/v1_1/transaction_handler.go (TransactionHandler.ViewCashoutTransaction)
GET    /v1_1/users/:guid/transactions/:transaction_guid/referral_cashback_transactions         --> /application/v1_1/transaction_handler.go (TransactionHandler.ViewReferralCashbackTransaction)
POST   /v1_1/users/:guid/transactions/cashout_transactions                                     --> /application/v1_1/cashout_transaction_handler.go (CashoutTransactionHandler.Create)
POST   /v1_1/users/:guid/edm/insufficient_funds                                                --> /application/v1_1/edm_handler.go (EdmHandler.InsufficientFunds)
GET    /v1_1/device/:device_uuid/users/:user_guid/notifications                                --> /application/v1_1/notification_handler.go (NotificationHandler.ViewNotificationForRegisteredUser)
```

### Development Workflow

1. Create new branch and naming it according to the reference below:
```
- Branch name start with `feature_` if you want to build new features. For example `feature_notification`
- Branch name start with `hotfix_` if you want to fix bug that effect production environment and need to solve it immediately.
- Branch name start with `bugfix_` if you want to fix less important bug and can push the fix later.
- Other than that, you can naming the branch with any word. For example `update_readme`.
```

2. Make any change to the current code and commit. Then push new branch to git repository.

3. After bitbucket detect changes in the repository, Bitbucket will trigger jenkins to build, test and output the result of the project.

4. If jenkins build success, create pull request to develop (Staging Environment).

![Jenkins Build Result Success](https://s3-ap-southeast-1.amazonaws.com/shoppermate/api_documentation_images/jenkins_build_result.png "Jenkins Build Result Success")

5. If jenkins build failed, check why it's failed and fix it. Then, repeat step 2. 

![Jenkins Build Result Failed](https://s3-ap-southeast-1.amazonaws.com/shoppermate/api_documentation_images/jenkins_build_result_failed.png "Jenkins Build Result Failed")

6. Bitbucket will trigger jenkins one more time to build, test and output the result of the project.

7. If jenkins build success, approved pull request.

8. Deploy to staging.

9. After everything okay, create pull request from develop branch to master branch. Approved the pull request and deploy to production.

### Code Description

#### Middleware

- Middleware file located in `middlewares` folder. Right now this API only have one middleware call `Auth (auth.go)`.

- You can find this code `version1.Use(middlewares.Auth(DB))` in routes file. This code responsible to initialize `auth middleware`.

- All endpoint inside above code require access token.

- That's middleware used to check if access token exist. If not exist, API will return error

- API also check if the token valid or not by checking device exist or not in database using `device_uuid` and `user_guid`.

#### Handler

- Handler is like Controller if you familiar with MVC. Handler will control what to do with the request and when to end the request or response with some data.

- All handlers reside in `application/{api_version}/` and the naming always end with `_handler.go`. For example `UserHandler (user_handler.go)`.

#### Services

##### Resource Services

- Every resources has their own service. For example, for `User Service (user_service.go)` is service for `User` resource.

- Service file will control application logics. All services resize in `application/{api_version}` and the naming always end with `_service.go`.

- For example API want to create new user. First thing, API will check if user already exist by checking user phone number in database.

- If exist, return an error. If not exist continue with checking if user register using referral code.

- If referral code not found, return an error. If referral code found, continue request with another application logic required during create user.

#### Repository

- Repository handle all task related to CRUD function. All repositories reside in `application/{api_version}` and the naming always end with `_repository.go`. One of the repository is `UserRepository (user_repository.go)`

#### Model

- Model represent the data for the resources. For example User resource. All models reside in `application/{api_version}` and the naming always end with `_model.go`. One of the model is `User (user_repository.go)`

- Model also used to specify database relationship. For example User has many shopping list. You can see code below to specify the relationship.

```
ShoppingLists []*ShoppingList `json:"shopping_lists,omitempty" gorm:"ForeignKey:UserGUID;AssociationForeignKey:GUID"`
```

- For the details reference refer to the link below:

`http://jinzhu.me/gorm/associations.html`

#### Validation

##### Validate Request Body
- To validate parameter in body, you must specify validation rules in `binding` files. For example `auth_binding.go`, `user_binding.go`.

```go
type CreateUser struct {
	FacebookID     string `form:"facebook_id" json:"facebook_id" binding:"omitempty,numeric"`
	Name           string `form:"name" json:"name" binding:"required"`
	Email          string `form:"email" json:"email" binding:"required,email"`
	PhoneNo        string `form:"phone_no" json:"phone_no" binding:"required,numeric,min=11,max=13"`
	ProfilePicture string `form:"profile_picture" json:"profile_picture" binding:"omitempty"`
	ReferralCode   string `form:"referral_code" json:"referral_code" binding:"omitempty,alphanum,max=8"`
	Debug          int    `form:"debug" json:"debug" binding:"omitempty"`
}
```

- Based on example above, refer to binding section to know if the parameter got validation rule or not.

- When API receives request, API will bind request data into the struct specify in handler.

##### Validate Query String

- To validate query string, you can use validation class `systems/validation.go`  in system package.

- Validation class is a wrapper for `go-playground/validator` package. Validation class used that package to validate input. You also can add new validation rule inside the class.

| Available Validation Rules | Example Syntax
| ------------- | -------------------- |
| - required     | required
| - uuid5      | uuid5
| - alpha | alpha
| - alphanum | alphanum
| - numeric | numeric
| - min | min=5
| - max | max=5
| - email | email
| - len | len=5
| - time | time
| - latitude | latitude
| - longitude | longitude

- For example, you can see `application/v11/deal_cashback_handler.go` in `ViewByShoppingList` function.

```go
// Specify Validation Rules
queryStringValidationRules := map[string]string {
		"page_number": "numeric",
		"page_limit":  "numeric",
}

// Validate input from query string
error := Validation.Validate(context.Request.URL.Query(), queryStringValidationRules)

// Return error if validation failed
if error != nil {
	context.JSON(http.StatusUnprocessableEntity, error)
	return
}
```

#### Validate File

- To validate file upload, you can use `FileValidation` class inside `services/filesystem/filesystem.go`

- API can validate file type and file size.

- For example, see function `UploadUserProfilePicture` in `application/v11/user_service.go`. That function responsible to validate and upload profile image to Amazon S3.

```go
// Validate if profile picture type is `jpg`, `jpeg`, `png` or `gif`
error := FileValidation.ValidateFileType([]string{"jpg", "jpeg", "png", "gif"}, profilePicture)

// Return error if profile picture type is different from above.
if error != nil {
	return nil, error
}

// Validate profile picture size. File size must specify in byte. Profile picture size must below than 1MB.
_, error = FileValidation.ValidateFileSize(profilePicture, 1000000, "profile_picture")

if error != nil {
	return nil, error
}

// Specify local upload path. API will upload to local first before upload to amazon S3. After finish uploading to Amazon S3, API will delete the local file.
localUploadPath := os.Getenv("GOPATH") + os.Getenv("STORAGE_PATH")

amazonS3UploadPath := "/profile_images/"

// Upload file to Amazon S3
uploadedFile, error := us.AmazonS3FileSystem.Upload(profilePicture, localUploadPath, amazonS3UploadPath)

if error != nil {
	return nil, error
}

return uploadedFile, nil
```

#### Request Lifecycle

- When API receive any request, API pass the request to the matched handler.

- Handler will bind and validate request data based on struct in binding file.

- Handler pass request data to Service. Service process the application logic.

- Service will use Repository if the application login require to handle task related to CRUD.

- Service will return back the data to handler and handler will output the result in JSON format.