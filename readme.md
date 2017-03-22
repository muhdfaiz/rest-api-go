## Shoppermate API

[![Build Status](http://188.166.227.158:8080/buildStatus/icon?job=Shoppermate API](http://188.166.227.158:8080/job/Shoppermate%20API/

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

#### Install Go Language 1.7.x
- Download Go Language 1.7.x
```
sudo apt-get update
sudo wget https://storage.go ogleapis.com/golang/go1.7.4.linux-amd64.tar.gz
```

- Extract Go Language 1.7.x
```
sudo tar -xvf go1.7.linux-amd64.tar.
sudo mv go /usr/local
```

#### Setup Go Environment.

Edit file `~/.profile` and include 3 environment variables below.

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

#### Setting Up Shoppermate API
- Allow go get to retrieve shoppermate API from private bitbucket repositories. Enter code below in command line.

```
git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
```

Here, we say use git@bitbucket.org any time you’d use https://bitbucket.org. This works for everything, not just go get. It just has the nice side effect of using your SSH key any time you run go get too.

- Verify .gitconfig file contain information like below
```
[url "git@bitbucket.org:"]
        insteadOf = https://bitbucket.org/
```

- Generate Your SSH Key. Paste code below in your command line and press enter until finish.
```
ssh-keygen
```

- Copy the ssh key out from the command below and add the SSH key in Bitbucket Repository.
```
cat ~/.ssh/id_rsa.pub
```

- Go to shoppermate-api project path. For example `~/golang/src/` and then install package dependencies using Glide.
Glide will install another package dependencies into `~/golang/src//vendor` folder
```
glide install
```

- Create new .env file and copy the content from .env.example file in root directory. Update all setting in .env file.

- Go to project root directory and run Shoppermate API.
```
go run api.go 
```

### Project Documentation

#### Library Used
- ORM: [GORM](http://jinzhu.me/gorm/
- Framework: [GIN](https://github.com/gin-gonic/gin
- Validator: [Go Playground Validator v8](https://github.com/go-playground/validator

#### Routes

- File Location

```
application/v1/routes.go 
application/v1_1/routes.go 
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