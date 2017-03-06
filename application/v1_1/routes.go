package v1_1

import (
	"os"

	"bitbucket.org/cliqers/shoppermate-api/middlewares"
	"bitbucket.org/cliqers/shoppermate-api/services/email"
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitializeObjectAndSetRoutesV1_1 will initialize object and set all routes across the API version 1.1
func InitializeObjectAndSetRoutesV1_1(router *gin.Engine, DB *gorm.DB) *gin.Engine {
	router.Use(func(context *gin.Context) {
		context.Set("DB", DB)
	})

	// Amazon S3 Config
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION_NAME")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	// Amazon S3 filesystem Objects
	fileSystem := &filesystem.FileSystem{}
	amazonS3FileSystem := fileSystem.Driver("amazonS3").(*filesystem.AmazonS3Upload)
	amazonS3FileSystem.AccessKey = accessKey
	amazonS3FileSystem.SecretKey = secretKey
	amazonS3FileSystem.Region = region
	amazonS3FileSystem.BucketName = bucketName

	// Device Objects
	deviceRepository := &DeviceRepository{DB: DB}
	deviceService := &DeviceService{DeviceRepository: deviceRepository}

	// Sms Objects
	smsHistoryRepository := &SmsHistoryRepository{DB: DB}
	smsHistoryService := &SmsHistoryService{SmsHistoryRepository: smsHistoryRepository}
	smsService := &SmsService{DB: DB, SmsHistoryRepository: smsHistoryRepository}

	// Auth Service Objects
	authService := &AuthService{SmsService: smsService, DeviceService: deviceService}

	// Facebook Service Objects
	facebookService := &facebook.FacebookService{}

	// Email Service Objects
	emailService := &email.EmailService{}

	// Occasion Objects
	occasionRepository := &OccasionRepository{DB: DB}
	occasionTransformer := &OccasionTransformer{}
	occasionService := &OccasionService{OccasionRepository: occasionRepository, OccasionTransformer: occasionTransformer}

	// Item Objects
	itemRepository := &ItemRepository{DB: DB}
	itemTransformer := &ItemTransformer{}
	itemService := &ItemService{ItemRepository: itemRepository}

	// Deal Cashback Repository Objects
	dealCashbackRepository := &DealCashbackRepository{DB: DB}

	// LocationService Objects
	locationService := &location.LocationService{}

	// Grocer Location Repository Objects
	grocerLocationRepository := &GrocerLocationRepository{DB: DB}

	// Grocer Location Service Objects
	grocerLocationService := &GrocerLocationService{GrocerLocationRepository: grocerLocationRepository}

	// Deal Repository Objects
	dealRepository := &DealRepository{DB: DB, GrocerLocationService: grocerLocationService}

	// Item SubCategory Objects
	itemSubCategoryRepository := &ItemSubCategoryRepository{DB: DB}
	itemSubCategoryService := &ItemSubCategoryService{ItemSubCategoryRepository: itemSubCategoryRepository, DealRepository: dealRepository}

	// Grocer Objects
	grocerRepository := &GrocerRepository{DB: DB}
	grocerService := &GrocerService{GrocerRepository: grocerRepository, DealRepository: dealRepository}

	// Item Category Objects
	itemCategoryRepository := &ItemCategoryRepository{DB: DB}
	itemCategoryTransformer := &ItemCategoryTransformer{}

	shoppingListItemRepository := &ShoppingListItemRepository{DB: DB}

	// Deal Transformer
	dealTransformer := &DealTransformer{}

	// Deal Service
	dealService := &DealService{DealRepository: dealRepository, DealTransformer: dealTransformer, LocationService: locationService,
		DealCashbackRepository: dealCashbackRepository, ItemRepository: itemRepository, ItemCategoryRepository: itemCategoryRepository,
		ItemSubCategoryRepository: itemSubCategoryRepository, GrocerService: grocerService, ShoppingListItemRepository: shoppingListItemRepository}

	itemCategoryService := &ItemCategoryService{ItemCategoryRepository: itemCategoryRepository, DealService: dealService,
		ItemCategoryTransformer: itemCategoryTransformer, GrocerRepository: grocerRepository, GrocerService: grocerService,
		DealRepository: dealRepository}

	// Default Shopping List Objects
	defaultShoppingListRepository := &DefaultShoppingListRepository{DB: DB}
	defaultShoppingListService := &DefaultShoppingListService{DefaultShoppingListRepository: defaultShoppingListRepository, DealService: dealService}

	// Default Shopping List Item Objects
	defaultShoppingListItemRepository := &DefaultShoppingListItemRepository{DB: DB}
	defaultShoppingListItemService := &DefaultShoppingListItemService{DefaultShoppingListItemRepository: defaultShoppingListItemRepository}

	// Generic Objects
	genericRepository := &GenericRepository{DB: DB}
	genericService := &GenericService{GenericRepository: genericRepository}
	genericTransformer := &GenericTransformer{}

	// Shopping List Item Service Object
	shoppingListItemService := &ShoppingListItemService{ShoppingListItemRepository: shoppingListItemRepository, ItemService: itemService,
		ItemCategoryService: itemCategoryService, ItemSubCategoryService: itemSubCategoryService, DealService: dealService, GenericService: genericService}

	// Shopping List Item Image Objects
	shoppingListItemImageRepository := &ShoppingListItemImageRepository{DB: DB}
	shoppingListItemImageService := &ShoppingListItemImageService{DB: DB, AmazonS3FileSystem: amazonS3FileSystem,
		ShoppingListItemImageRepository: shoppingListItemImageRepository, ShoppingListItemService: shoppingListItemService}

	// Shopping List Objects
	shoppingListRepository := &ShoppingListRepository{DB: DB}
	shoppingListService := &ShoppingListService{ShoppingListRepository: shoppingListRepository, OccasionService: occasionService,
		DefaultShoppingListService: defaultShoppingListService, DefaultShoppingListItemService: defaultShoppingListItemService,
		ShoppingListItemService: shoppingListItemService, ShoppingListItemImageService: shoppingListItemImageService}

	// Deal Cashback Transformer Objects
	dealCashbackTransformer := &DealCashbackTransformer{}

	// Deal Cashback Service Objects
	dealCashbackService := &DealCashbackService{DealCashbackRepository: dealCashbackRepository,
		ShoppingListService: shoppingListService, ShoppingListItemService: shoppingListItemService, DealRepository: dealRepository}

	// Transaction Status Objects
	transactionStatusRepository := &TransactionStatusRepository{DB: DB}
	transactionStatusService := &TransactionStatusService{TransactionStatusRepository: transactionStatusRepository}

	// Transaction Type Objects
	transactionTypeRepository := &TransactionTypeRepository{DB: DB}
	transactionTypeService := &TransactionTypeService{TransactionTypeRepository: transactionTypeRepository}

	// Transaction Objects
	transactionRepository := &TransactionRepository{DB: DB, TransactionStatusRepository: transactionStatusRepository}
	transactionTransformer := &TransactionTransformer{}
	transactionService := &TransactionService{TransactionRepository: transactionRepository, TransactionTransformer: transactionTransformer,
		DealCashbackService: dealCashbackService, ShoppingListService: shoppingListService, DealService: dealService,
		TransactionTypeService: transactionTypeService, TransactionStatusService: transactionStatusService}

	// Referral Cashback Transaction Objects
	referralCashbackTransactionRepository := &ReferralCashbackTransactionRepository{DB: DB}
	referralCashbackTransactionService := &ReferralCashbackTransactionService{ReferralCashbackTransactionRepository: referralCashbackTransactionRepository}

	// User Objects Objects
	userRepository := &UserRepository{DB: DB}
	userService := &UserService{UserRepository: userRepository, TransactionService: transactionService, TransactionStatusService: transactionStatusService,
		TransactionTypeService: transactionTypeService, DealCashbackService: dealCashbackService, FacebookService: facebookService, SmsService: smsService,
		DeviceService: deviceService, ReferralCashbackTransactionService: referralCashbackTransactionService, AmazonS3FileSystem: amazonS3FileSystem,
		SmsHistoryService: smsHistoryService, EmailService: emailService}

	//Deal Cashback Transaction
	dealCashbackTransactionRepository := &DealCashbackTransactionRepository{DB: DB}
	dealCashbackTransactionService := &DealCashbackTransactionService{AmazonS3FileSystem: amazonS3FileSystem,
		DealCashbackRepository: dealCashbackRepository, DealCashbackTransactionRepository: dealCashbackTransactionRepository,
		DealRepository: dealRepository, TransactionRepository: transactionRepository, ShoppingListItemRepository: shoppingListItemRepository,
		TransactionStatusService: transactionStatusService, TransactionTypeService: transactionTypeService}

	// Event Objects
	eventRepository := &EventRepository{DB: DB}
	eventService := &EventService{EventRepository: eventRepository, DealCashbackRepository: dealCashbackRepository, DealService: dealService}

	//Cashout Transaction Object
	cashoutTransactionRepository := &CashoutTransactionRepository{DB: DB}
	cashoutTransactionService := &CashoutTransactionService{CashoutTransactionRepository: cashoutTransactionRepository,
		TransactionService: transactionService, TransactionStatusService: transactionStatusService,
		TransactionTypeService: transactionTypeService, UserRepository: userRepository, EmailService: emailService}

	// Setting Objects
	settingRepository := &SettingRepository{DB: DB}
	settingService := &SettingService{SettingRepository: settingRepository}

	// Featured Deal Objects
	featuredDealRepository := &FeaturedDealRepository{DB: DB}

	// EDM History Objects
	edmHistoryRepository := &EdmHistoryRepository{DB: DB}

	// EDM Objects
	edmService := &EdmService{EmailService: emailService, DealService: dealService, FeaturedDealRepository: featuredDealRepository,
		EdmHistoryRepository: edmHistoryRepository}

	// Notification Objects
	notificationRepository := &NotificationRepository{DB: DB}
	notificationService := &NotificationService{NotificationRepository: notificationRepository}

	// Sms Handler
	smsHandler := &SmsHandler{UserRepository: userRepository, SmsService: smsService,
		SmsHistoryService: smsHistoryService, DeviceService: deviceService, UserService: userService}

	// User Handler
	userHandler := &UserHandler{UserService: userService, SettingService: settingService}

	// Device Handler
	deviceHandler := &DeviceHandler{DeviceService: deviceService, UserService: userService}

	// Shopping List Handler
	shoppingListHandler := &ShoppingListHandler{ShoppingListService: shoppingListService, ShoppingListItemImageService: shoppingListItemImageService}

	// Auth Handler
	authHandler := &AuthHandler{AuthService: authService, UserService: userService}

	// Occasion Handler
	occasionHandler := &OccasionHandler{OccasionService: occasionService}

	// Item Handler
	itemHandler := &ItemHandler{ItemRepository: itemRepository, ItemTransformer: itemTransformer}

	// Item Category Handler
	itemCategoryHandler := &ItemCategoryHandler{ItemCategoryService: itemCategoryService}

	// Shopping List Item Handler
	shoppingListItemHandler := &ShoppingListItemHandler{ShoppingListService: shoppingListService, ShoppingListItemService: shoppingListItemService,
		ShoppingListItemImageService: shoppingListItemImageService}

	// Shopping List Item Image Handler
	shoppingListItemImageHandler := &ShoppingListItemImageHandler{ShoppingListItemImageService: shoppingListItemImageService}

	// Deal Cashback Handler
	dealCashbackHandler := &DealCashbackHandler{ShoppingListRepository: shoppingListRepository, DealService: dealService, DealCashbackService: dealCashbackService,
		DealCashbackTransformer: dealCashbackTransformer}

	// Deal Handler
	dealHandler := &DealHandler{DealService: dealService, DealTransformer: dealTransformer, ItemCategoryService: itemCategoryService,
		DealCashbackService: dealCashbackService, UserRepository: userRepository, ItemSubCategoryRepository: itemSubCategoryRepository}

	// Event Handler
	eventHandler := &EventHandler{EventService: eventService}

	// Deal Cashback Transaction Handler
	dealCashbackTransactionHandler := &DealCashbackTransactionHandler{DealCashbackTransactionService: dealCashbackTransactionService, TransactionService: transactionService}

	// Transaction Handler
	transactionHandler := &TransactionHandler{TransactionService: transactionService}

	// Cashout Transaction Handler
	cashoutTransactionHandler := &CashoutTransactionHandler{CashoutTransactionService: cashoutTransactionService, TransactionService: transactionService}

	// Default Shopping List Handler
	defaultShoppingListHandler := &DefaultShoppingListHandler{DefaultShoppingListService: defaultShoppingListService}

	// Grocer Handler
	grocerHandler := &GrocerHandler{GrocerService: grocerService}

	// Generic Handler
	genericHandler := &GenericHandler{GenericService: genericService, GenericTransformer: genericTransformer}

	// Setting Handler
	settingHandler := &SettingHandler{SettingService: settingService}

	// EDM Handler
	edmHandler := &EdmHandler{EdmService: edmService}

	// Notification Handler
	notificationHandler := &NotificationHandler{NotificationService: notificationService, DeviceService: deviceService}

	//  Routes
	version1_1 := router.Group("/v1_1")
	{
		// Public Routes
		// Device Routes
		version1_1.POST("/devices", deviceHandler.Create)
		version1_1.PATCH("/devices/:uuid", deviceHandler.Update)

		// User Routes
		version1_1.POST("/users", userHandler.Create)

		// SMS Routes
		version1_1.POST("/sms", smsHandler.Send)
		version1_1.POST("/sms/verifications", smsHandler.Verify)

		// Authentication Routes
		version1_1.POST("/auth/login/phone", authHandler.LoginViaPhone)
		version1_1.POST("/auth/login/facebook", authHandler.LoginViaFacebook)

		// Occasion Routes
		version1_1.GET("/shopping_lists/occasions", occasionHandler.Index)

		// Shopping List Item Routes
		version1_1.GET("/shopping_lists/items", itemHandler.Index)

		// Shopping List Item Categories Routes
		version1_1.GET("/shopping_lists/items/categories", itemCategoryHandler.ViewAll)

		// Generic Category Routes
		version1_1.GET("generics", genericHandler.ViewAll)

		// Deal Route
		version1_1.GET("deals", dealHandler.ViewAllForGuestUser)

		// Default Shopping List Route
		version1_1.GET("/shopping_list_samples", defaultShoppingListHandler.ViewAll)

		// Setting  Route
		version1_1.GET("/settings", settingHandler.ViewAll)

		// Notification Routes
		version1_1.GET("/device/:device_uuid/notifications", notificationHandler.ViewNotificationForGuest)

		// Protected Routes
		version1_1.Use(middlewares.Auth(DB))
		{
			// User Routes
			version1_1.PATCH("/users/:guid", userHandler.Update)
			version1_1.GET("/users/:guid", userHandler.View)

			// Device Routes
			version1_1.DELETE("/devices/:uuid", deviceHandler.Delete)

			// Authentication Routes
			version1_1.GET("/auth/refresh", authHandler.Refresh)
			version1_1.GET("/auth/logout", authHandler.Logout)

			// User Shopping List Routes
			version1_1.GET("users/:guid/shopping_lists", shoppingListHandler.View)
			version1_1.POST("users/:guid/shopping_lists", shoppingListHandler.Create)
			version1_1.PATCH("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Update)
			version1_1.DELETE("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Delete)

			// User Shopping List Item Routes
			version1_1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.View)
			version1_1.GET("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.ViewAll)
			version1_1.POST("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.Create)
			version1_1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Update)
			version1_1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.UpdateAll)
			version1_1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Delete)
			version1_1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.DeleteAll)

			// User Shopping List Item Image Routes
			version1_1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guid", shoppingListItemImageHandler.View)
			version1_1.POST("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images", shoppingListItemImageHandler.Create)
			version1_1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guids", shoppingListItemImageHandler.Delete)

			// Deal Routes
			version1_1.GET("users/:guid/deals", dealHandler.ViewAllForRegisteredUser)
			version1_1.GET("deals/:deal_guid", dealHandler.View)
			version1_1.GET("users/:guid/deals/categories", dealHandler.ViewAndGroupByCategory)
			version1_1.GET("users/:guid/deals/categories/:category_guid", dealHandler.ViewByCategory)
			version1_1.GET("users/:guid/deals/categories/:category_guid/subcategories", dealHandler.ViewByCategoryAndGroupBySubCategory)
			version1_1.GET("users/:guid/deals/subcategories/:subcategory_guid", dealHandler.ViewBySubCategory)
			version1_1.GET("users/:guid/deals/grocers", grocerHandler.GetAllGrocersThatContainDeals)
			version1_1.GET("users/:guid/deals/grocers/:grocer_guid/categories", itemCategoryHandler.ViewGrocerCategoriesThoseHaveDealsIncludingDeals)
			version1_1.GET("users/:guid/deals/grocers/:grocer_guid/categories/:category_guid", dealHandler.ViewByGrocerAndCategory)

			// Feature Deal Routes
			version1_1.GET("users/:guid/featured_deals", eventHandler.ViewAll)

			// Deal Cashback Routes
			version1_1.POST("users/:guid/deal_cashbacks", dealCashbackHandler.Create)
			version1_1.GET("users/:guid/deal_cashbacks/shopping_lists/:shopping_list_guid", dealCashbackHandler.ViewByShoppingList)
			version1_1.GET("users/:guid/deal_cashbacks", dealCashbackHandler.ViewByUserAndGroupByShoppingList)
			version1_1.GET("users/:guid/deal_cashbacks/deals/:deal_guid", dealCashbackHandler.ViewByUserAndDealGroupByShoppingList)

			// Deal Cashback Transaction Routes
			version1_1.POST("users/:guid/transactions/deal_cashback_transactions", dealCashbackTransactionHandler.Create)

			// Transaction Routes Routes
			version1_1.GET("users/:guid/transactions", transactionHandler.ViewUserTransactions)
			version1_1.GET("users/:guid/transactions/:transaction_guid/deal_cashback_transactions", transactionHandler.ViewDealCashbackTransaction)
			version1_1.GET("users/:guid/transactions/:transaction_guid/cashout_transactions", transactionHandler.ViewCashoutTransaction)
			version1_1.GET("users/:guid/transactions/:transaction_guid/referral_cashback_transactions", transactionHandler.ViewReferralCashbackTransaction)

			// Cashout Transaction Routes
			version1_1.POST("users/:guid/transactions/cashout_transactions", cashoutTransactionHandler.Create)

			// EDM  Routes
			version1_1.POST("/users/:guid/edm/insufficient_funds", edmHandler.InsufficientFunds)

			// Notification Routes
			version1_1.GET("/device/:device_uuid/users/:user_guid/notifications", notificationHandler.ViewNotificationForRegisteredUser)
		}
	}

	return router
}
