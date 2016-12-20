package application

import (
	"os"

	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/middlewares"
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitializeObjectAndSetRoutes will initialize object and set all routes across the API
func InitializeObjectAndSetRoutes(router *gin.Engine, DB *gorm.DB) *gin.Engine {
	// Amazon S3 Config
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION_NAME")
	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")

	// Amazon S3 filesystem
	fileSystem := &filesystem.FileSystem{}
	amazonS3FileSystem := fileSystem.Driver("amazonS3").(*filesystem.AmazonS3Upload)
	amazonS3FileSystem.AccessKey = accessKey
	amazonS3FileSystem.SecretKey = secretKey
	amazonS3FileSystem.Region = region
	amazonS3FileSystem.BucketName = bucketName

	// Referral Cashback Objects
	referralCashbackRepository := &v1.ReferralCashbackRepository{DB: DB}

	// User Objects
	userRepository := &v1.UserRepository{DB: DB}
	userService := &v1.UserService{DB: DB, AmazonS3FileSystem: amazonS3FileSystem, UserRepository: userRepository,
		ReferralCashbackRepository: referralCashbackRepository}

	// Device Objects
	deviceRepository := &v1.DeviceRepository{DB: DB}
	deviceService := &v1.DeviceService{DeviceRepository: deviceRepository, UserService: userService}

	// Sms Objects
	smsHistoryRepository := &v1.SmsHistoryRepository{DB: DB}
	smsService := &v1.SmsService{DB: DB, SmsHistoryRepository: smsHistoryRepository}

	// Auth Service
	authService := &v1.AuthService{UserService: userService, SmsService: smsService, DeviceService: deviceService}

	// Facebook Service
	facebookService := &facebook.FacebookService{
		AppID:     os.Getenv("FACEBOOK_APP_ID"),
		AppSecret: os.Getenv("FACEBOOK_APP_SECRET"),
	}

	// Occasion Objects
	occasionRepository := &v1.OccasionRepository{DB: DB}
	occasionTransformer := &v1.OccasionTransformer{}
	occasionService := &v1.OccasionService{OccasionRepository: occasionRepository, OccasionTransformer: occasionTransformer}

	// Item Objects
	itemRepository := &v1.ItemRepository{DB: DB}
	itemTransformer := &v1.ItemTransformer{}
	itemService := &v1.ItemService{ItemRepository: itemRepository}

	// Deal Cashback Repository
	dealCashbackRepository := &v1.DealCashbackRepository{DB: DB}

	// LocationService
	locationService := &location.LocationService{}

	// Grocer Location Repository
	grocerLocationRepository := &v1.GrocerLocationRepository{DB: DB}

	// Grocer Location Service
	grocerLocationService := &v1.GrocerLocationService{GrocerLocationRepository: grocerLocationRepository}

	// Deal Repository
	dealRepository := &v1.DealRepository{DB: DB, GrocerLocationService: grocerLocationService}

	// Item SubCategory Objects
	itemSubCategoryRepository := &v1.ItemSubCategoryRepository{DB: DB}
	itemSubCategoryService := &v1.ItemSubCategoryService{ItemSubCategoryRepository: itemSubCategoryRepository, DealRepository: dealRepository}

	// Grocer Objects
	grocerRepository := &v1.GrocerRepository{DB: DB}
	grocerService := &v1.GrocerService{GrocerRepository: grocerRepository, DealRepository: dealRepository}

	// Item Category Objects
	itemCategoryRepository := &v1.ItemCategoryRepository{DB: DB}
	itemCategoryTransformer := &v1.ItemCategoryTransformer{}

	shoppingListItemRepository := &v1.ShoppingListItemRepository{DB: DB}

	// Deal Transformer
	dealTransformer := &v1.DealTransformer{}

	// Deal Service
	dealService := &v1.DealService{DealRepository: dealRepository, DealTransformer: dealTransformer, LocationService: locationService,
		DealCashbackRepository: dealCashbackRepository, ItemRepository: itemRepository, ItemCategoryRepository: itemCategoryRepository,
		ItemSubCategoryRepository: itemSubCategoryRepository, GrocerService: grocerService, ShoppingListItemRepository: shoppingListItemRepository}

	itemCategoryService := &v1.ItemCategoryService{ItemCategoryRepository: itemCategoryRepository, DealService: dealService,
		ItemCategoryTransformer: itemCategoryTransformer, GrocerRepository: grocerRepository, GrocerService: grocerService,
		DealRepository: dealRepository}

	// Default Shopping List Objects
	defaultShoppingListRepository := &v1.DefaultShoppingListRepository{DB: DB}
	defaultShoppingListService := &v1.DefaultShoppingListService{DefaultShoppingListRepository: defaultShoppingListRepository, DealService: dealService}

	defaultShoppingListItemRepository := &v1.DefaultShoppingListItemRepository{DB: DB}
	defaultShoppingListItemService := &v1.DefaultShoppingListItemService{DefaultShoppingListItemRepository: defaultShoppingListItemRepository}

	// Shopping List Item Repository
	shoppingListItemService := &v1.ShoppingListItemService{ShoppingListItemRepository: shoppingListItemRepository, ItemService: itemService,
		ItemCategoryService: itemCategoryService, ItemSubCategoryService: itemSubCategoryService, DealService: dealService}

	// Shopping List Item Image Objects
	shoppingListItemImageRepository := &v1.ShoppingListItemImageRepository{DB: DB}
	shoppingListItemImageService := &v1.ShoppingListItemImageService{DB: DB, AmazonS3FileSystem: amazonS3FileSystem,
		ShoppingListItemImageRepository: shoppingListItemImageRepository, ShoppingListItemService: shoppingListItemService}

	// Shopping List Objects
	shoppingListRepository := &v1.ShoppingListRepository{DB: DB}
	shoppingListService := &v1.ShoppingListService{ShoppingListRepository: shoppingListRepository, OccasionService: occasionService,
		DefaultShoppingListService: defaultShoppingListService, DefaultShoppingListItemService: defaultShoppingListItemService,
		ShoppingListItemService: shoppingListItemService, ShoppingListItemImageService: shoppingListItemImageService}

	// Deal Cashback Transformer
	dealCashbackTransformer := &v1.DealCashbackTransformer{}

	// Deal Cashback Service
	dealCashbackService := &v1.DealCashbackService{DealCashbackRepository: dealCashbackRepository,
		ShoppingListItemService: shoppingListItemService, DealRepository: dealRepository}

	// Transaction Status
	transactionStatusRepository := &v1.TransactionStatusRepository{DB: DB}
	transactionStatusService := &v1.TransactionStatusService{TransactionStatusRepository: transactionStatusRepository}

	// Transaction Type
	transactionTypeRepository := &v1.TransactionTypeRepository{DB: DB}
	transactionTypeService := &v1.TransactionTypeService{TransactionTypeRepository: transactionTypeRepository}

	// Transaction
	transactionRepository := &v1.TransactionRepository{DB: DB, TransactionStatusRepository: transactionStatusRepository}
	transactionTransformer := &v1.TransactionTransformer{}
	transactionService := &v1.TransactionService{TransactionRepository: transactionRepository, TransactionTransformer: transactionTransformer,
		DealCashbackService: dealCashbackService, ShoppingListService: shoppingListService, DealService: dealService,
		TransactionTypeService: transactionTypeService, TransactionStatusService: transactionStatusService}

	//Deal Cashback Transaction
	dealCashbackTransactionRepository := &v1.DealCashbackTransactionRepository{DB: DB}
	dealCashbackTransactionService := &v1.DealCashbackTransactionService{AmazonS3FileSystem: amazonS3FileSystem,
		DealCashbackRepository: dealCashbackRepository, DealCashbackTransactionRepository: dealCashbackTransactionRepository,
		DealRepository: dealRepository, TransactionTypeRepository: transactionTypeRepository, TransactionRepository: transactionRepository}

	// Generic Objects
	genericRepository := &v1.GenericRepository{DB: DB}
	genericService := &v1.GenericService{GenericRepository: genericRepository}
	genericTransformer := &v1.GenericTransformer{}

	// Event Objects
	eventRepository := &v1.EventRepository{DB: DB}
	eventService := &v1.EventService{EventRepository: eventRepository, DealCashbackRepository: dealCashbackRepository, DealService: dealService}

	//Cashout Transaction Object
	cashoutTransactionRepository := &v1.CashoutTransactionRepository{DB: DB}
	cashoutTransactionService := &v1.CashoutTransactionService{CashoutTransactionRepository: cashoutTransactionRepository,
		TransactionService: transactionService, UserRepository: userRepository, TransactionTypeRepository: transactionTypeRepository}

	// Sms Handler
	smsHandler := v1.SmsHandler{UserRepository: userRepository, SmsService: smsService,
		SmsHistoryRepository: smsHistoryRepository, DeviceService: deviceService}

	// User Handler
	userHandler := v1.UserHandler{UserRepository: userRepository, UserService: userService,
		DeviceService: deviceService, ReferralCashbackRepository: referralCashbackRepository, SmsService: smsService,
		FacebookService: facebookService, TransactionService: transactionService, DealCashbackService: dealCashbackService}

	// Device Handler
	deviceHandler := v1.DeviceHandler{DeviceService: deviceService}

	// Shopping List Handler
	shoppingListHandler := v1.ShoppingListHandler{ShoppingListService: shoppingListService, ShoppingListItemImageService: shoppingListItemImageService}

	// Auth Handler
	authHandler := v1.AuthHandler{AuthService: authService}

	// Occasion Handler
	occasionHandler := v1.OccasionHandler{OccasionService: occasionService}

	// Item Handler
	itemHandler := v1.ItemHandler{ItemRepository: itemRepository, ItemTransformer: itemTransformer}

	// Item Category Handler
	itemCategoryHandler := v1.ItemCategoryHandler{ItemCategoryService: itemCategoryService}

	// Shopping List Item Handler
	shoppingListItemHandler := v1.ShoppingListItemHandler{ShoppingListService: shoppingListService, ShoppingListItemService: shoppingListItemService,
		ShoppingListItemImageService: shoppingListItemImageService}

	// Shopping List Item Image Handler
	shoppingListItemImageHandler := v1.ShoppingListItemImageHandler{ShoppingListItemImageService: shoppingListItemImageService}

	// Deal Cashback Handler
	dealCashbackHandler := v1.DealCashbackHandler{ShoppingListRepository: shoppingListRepository, DealCashbackService: dealCashbackService,
		DealCashbackTransformer: dealCashbackTransformer}

	// Deal Handler
	dealHandler := v1.DealHandler{DealService: dealService, DealTransformer: dealTransformer, ItemCategoryService: itemCategoryService,
		DealCashbackService: dealCashbackService, UserRepository: userRepository, ItemSubCategoryRepository: itemSubCategoryRepository}

	// Event Handler
	eventHandler := v1.EventHandler{EventService: eventService}

	// Deal Cashback Transaction Handler
	dealCashbackTransactionHandler := v1.DealCashbackTransactionHandler{DealCashbackTransactionService: dealCashbackTransactionService}

	transactionHandler := v1.TransactionHandler{TransactionService: transactionService}

	cashoutTransactionHandler := v1.CashoutTransactionHandler{CashoutTransactionService: cashoutTransactionService}

	defaultShoppingListHandler := v1.DefaultShoppingListHandler{DefaultShoppingListService: defaultShoppingListService}

	grocerHandler := v1.GrocerHandler{GrocerService: grocerService}

	genericHandler := v1.GenericHandler{GenericService: genericService, GenericTransformer: genericTransformer}

	// V1 Routes
	version1 := router.Group("/v1")
	{
		// Public Routes
		// Device Routes
		version1.POST("/devices", deviceHandler.Create)
		version1.PATCH("/devices/:uuid", deviceHandler.Update)

		// User Routes
		version1.POST("/users", userHandler.Create)

		// SMS Routes
		version1.POST("/sms", smsHandler.Send)
		version1.POST("/sms/verifications", smsHandler.Verify)

		// Authentication Routes
		version1.POST("/auth/login/phone", authHandler.LoginViaPhone)
		version1.POST("/auth/login/facebook", authHandler.LoginViaFacebook)

		// Occasion Routes
		version1.GET("/shopping_lists/occasions", occasionHandler.Index)

		// Shopping List Item Routes
		version1.GET("/shopping_lists/items", itemHandler.Index)

		// Shopping List Item Categories Routes
		version1.GET("/shopping_lists/items/categories", itemCategoryHandler.ViewAll)

		// Generic Category Routes
		version1.GET("generics", genericHandler.ViewAll)

		// Deal Route
		version1.GET("deals", dealHandler.ViewAllForGuestUser)

		// Default Shopping List Route
		version1.GET("/shopping_list_samples", defaultShoppingListHandler.ViewAll)

		// Protected Routes
		version1.Use(middlewares.Auth(DB))
		{
			// User Routes
			version1.PATCH("/users/:guid", userHandler.Update)
			version1.GET("/users/:guid", userHandler.View)

			// Device Routes
			version1.DELETE("/devices/:uuid", deviceHandler.Delete)

			// Authentication Routes
			version1.GET("/auth/refresh", authHandler.Refresh)
			version1.GET("/auth/logout", authHandler.Logout)

			// User Shopping List Routes
			version1.GET("users/:guid/shopping_lists", shoppingListHandler.View)
			version1.POST("users/:guid/shopping_lists", shoppingListHandler.Create)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Update)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Delete)

			// User Shopping List Item Routes
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.View)
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.ViewAll)
			version1.POST("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.Create)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Update)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.UpdateAll)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Delete)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.DeleteAll)

			// User Shopping List Item Image Routes
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guid", shoppingListItemImageHandler.View)
			version1.POST("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images", shoppingListItemImageHandler.Create)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guids", shoppingListItemImageHandler.Delete)

			// Deal Handler
			version1.GET("users/:guid/deals", dealHandler.ViewAllForRegisteredUser)
			version1.GET("deals/:deal_guid", dealHandler.View)
			version1.GET("users/:guid/deals/categories", dealHandler.ViewAndGroupByCategory)
			version1.GET("users/:guid/deals/categories/:category_guid", dealHandler.ViewByCategory)
			version1.GET("users/:guid/deals/categories/:category_guid/subcategories", dealHandler.ViewByCategoryAndGroupBySubCategory)
			version1.GET("users/:guid/deals/subcategories/:subcategory_guid", dealHandler.ViewBySubCategory)
			version1.GET("users/:guid/deals/grocers", grocerHandler.GetAllGrocersThatContainDeals)
			version1.GET("users/:guid/deals/grocers/:grocer_guid/categories", itemCategoryHandler.ViewGrocerCategoriesThoseHaveDealsIncludingDeals)
			version1.GET("users/:guid/deals/grocers/:grocer_guid/categories/:category_guid", dealHandler.ViewByGrocerAndCategory)

			// Feature Deal (In Carousel) Handler
			version1.GET("users/:guid/featured_deals", eventHandler.ViewAll)

			// Deal Cashback Handler
			version1.POST("users/:guid/deal_cashbacks", dealCashbackHandler.Create)
			version1.GET("users/:guid/deal_cashbacks/shopping_lists/:shopping_list_guid", dealCashbackHandler.ViewByShoppingList)

			// Deal Cashback Transaction
			version1.POST("users/:guid/transactions/deal_cashback_transactions", dealCashbackTransactionHandler.Create)

			// Transaction Routes
			version1.GET("users/:guid/transactions", transactionHandler.ViewUserTransactions)
			version1.GET("users/:guid/transactions/:transaction_guid/deal_cashback_transactions", transactionHandler.ViewDealCashbackTransaction)
			version1.GET("users/:guid/transactions/:transaction_guid/cashout_transactions", transactionHandler.ViewCashoutTransaction)

			// Cashout Transaction
			version1.POST("users/:guid/transactions/cashout_transactions", cashoutTransactionHandler.Create)

		}
	}

	return router
}
