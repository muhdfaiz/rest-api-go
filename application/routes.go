package application

import (
	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/middlewares"
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/services/filesystem"
	"bitbucket.org/cliqers/shoppermate-api/services/location"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

// InitializeObjectAndSetRoutes will initialize object and set all routes across the API
func InitializeObjectAndSetRoutes(router *gin.Engine) *gin.Engine {
	Database := &systems.Database{}
	DB := Database.Connect("production")

	router.Use(func(c *gin.Context) {
		DB := Database.Connect("production")
		c.Set("DB", DB)
		c.Next()
	})

	// Amazon S3 Config
	Config := &systems.Configs{}
	accessKey := Config.Get("app.yaml", "aws_access_key_id", "")
	secretKey := Config.Get("app.yaml", "aws_secret_access_key", "")
	region := Config.Get("app.yaml", "aws_region_name", "")
	bucketName := Config.Get("app.yaml", "aws_bucket_name", "")

	// Amazon S3 filesystem
	fileSystem := &filesystem.FileSystem{}
	amazonS3FileSystem := fileSystem.Driver("amazonS3").(*filesystem.AmazonS3Upload)
	amazonS3FileSystem.AccessKey = accessKey
	amazonS3FileSystem.SecretKey = secretKey
	amazonS3FileSystem.Region = region
	amazonS3FileSystem.BucketName = bucketName

	// User Objects
	userRepository := &v1.UserRepository{DB: DB}
	userFactory := &v1.UserFactory{DB: DB}
	userService := &v1.UserService{DB: DB, AmazonS3FileSystem: amazonS3FileSystem}

	// Device Objects
	deviceRepository := &v1.DeviceRepository{DB: DB}
	deviceFactory := &v1.DeviceFactory{DB: DB}

	// Sms Objects
	smsHistoryRepository := &v1.SmsHistoryRepository{DB: DB}
	smsService := &v1.SmsService{DB: DB}

	// Referral Cashback Objects
	referralCashbackRepository := &v1.ReferralCashbackRepository{DB: DB}

	// Facebook Service
	facebookService := &facebook.FacebookService{
		AppID:     Config.Get("app.yaml", "facebook_app_id", ""),
		AppSecret: Config.Get("app.yaml", "facebook_app_secret", ""),
	}

	// Occasion Objects
	occasionRepostory := &v1.OccasionRepository{DB: DB}
	occasionTransformer := &v1.OccasionTransformer{}

	// Item Objects
	itemRepository := &v1.ItemRepository{DB: DB}
	itemTransformer := &v1.ItemTransformer{}

	// Item Category Objects
	itemCategoryRepository := &v1.ItemCategoryRepository{DB: DB}
	itemCategoryTransformer := &v1.ItemCategoryTransformer{}
	itemCategoryService := &v1.ItemCategoryService{ItemCategoryRepository: itemCategoryRepository,
		ItemCategoryTransformer: itemCategoryTransformer}

	// Item SubCategory Objects
	itemSubCategoryRepository := &v1.ItemSubCategoryRepository{DB: DB}

	// Shopping List Objects
	shoppingListFactory := &v1.ShoppingListFactory{DB: DB}
	shoppingListRepository := &v1.ShoppingListRepository{DB: DB}

	// Shopping List Item Image Objects
	shoppingListItemImageService := &v1.ShoppingListItemImageService{DB: DB, AmazonS3FileSystem: amazonS3FileSystem}
	shoppingListItemImageFactory := &v1.ShoppingListItemImageFactory{DB: DB, ShoppingListItemImageService: shoppingListItemImageService}
	shoppingListItemImageRepository := &v1.ShoppingListItemImageRepository{DB: DB}

	// LocationService
	locationService := &location.LocationService{}

	dealCashbackFactory := &v1.DealCashbackFactory{DB: DB}

	// Grocer Location Repository
	grocerLocationRepository := &v1.GrocerLocationRepository{DB: DB}

	// Grocer Location Service
	grocerLocationService := &v1.GrocerLocationService{GrocerLocationRepository: grocerLocationRepository}

	// Deal Repository
	dealRepository := &v1.DealRepository{DB: DB, GrocerLocationService: grocerLocationService}

	// Shopping List Item Factory
	shoppingListItemFactory := &v1.ShoppingListItemFactory{DB: DB, ItemRepository: itemRepository, ShoppingListItemImageFactory: shoppingListItemImageFactory,
		ShoppingListItemImageRepository: shoppingListItemImageRepository, DealRepository: dealRepository, ItemCategoryRepository: itemCategoryRepository,
		ItemSubCategoryRepository: itemSubCategoryRepository}

	// Deal Cashback Repository
	dealCashbackRepository := &v1.DealCashbackRepository{DB: DB}

	// Deal Service
	dealService := &v1.DealService{DealRepository: dealRepository, LocationService: locationService, DealCashbackFactory: dealCashbackFactory,
		ShoppingListItemFactory: shoppingListItemFactory, DealCashbackRepository: dealCashbackRepository, ItemRepository: itemRepository,
		ItemCategoryService: itemCategoryService, ItemSubCategoryRepository: itemSubCategoryRepository}

	// Deal Transformer
	dealTransformer := &v1.DealTransformer{}

	// Shopping List Item Repository
	shoppingListItemRepository := &v1.ShoppingListItemRepository{DB: DB, DealService: dealService}

	// Deal Cashback Service
	dealCashbackService := &v1.DealCashbackService{DealCashbackRepository: dealCashbackRepository, DealRepository: dealRepository,
		DealCashbackFactory: dealCashbackFactory, ShoppingListItemFactory: shoppingListItemFactory}

	// Grocer Objects
	grocerRepository := &v1.GrocerRepository{DB: DB}
	grocerTransformer := &v1.GrocerTransformer{}

	// Event Objects
	eventRepository := &v1.EventRepository{DB: DB}
	eventService := &v1.EventService{EventRepository: eventRepository, DealCashbackRepository: dealCashbackRepository}

	// Sms Handler
	smsHandler := v1.SmsHandler{UserRepository: userRepository, UserFactory: userFactory, SmsService: smsService,
		SmsHistoryRepository: smsHistoryRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}

	// User Handler
	userHandler := v1.UserHandler{UserRepository: userRepository, UserService: userService, UserFactory: userFactory, DeviceFactory: deviceFactory,
		ReferralCashbackRepository: referralCashbackRepository, SmsService: smsService, FacebookService: facebookService}

	// Device Handler
	deviceHandler := v1.DeviceHandler{UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}

	// Shopping List Handler
	shoppingListHandler := v1.ShoppingListHandler{UserRepository: userRepository, OccasionRepository: occasionRepostory,
		ShoppingListFactory: shoppingListFactory, ShoppingListRepository: shoppingListRepository, ShoppingListItemFactory: shoppingListItemFactory}

	// Auth Handler
	authHandler := v1.AuthHandler{UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory,
		SmsService: smsService}

	// Occasion Handler
	occasionHandler := v1.OccasionHandler{OccasionRepository: occasionRepostory, OccasionTransformer: occasionTransformer}

	// Item Handler
	itemHandler := v1.ItemHandler{ItemRepository: itemRepository, ItemTransformer: itemTransformer}

	// Item Category Handler
	itemCategoryHandler := v1.ItemCategoryHandler{ItemCategoryService: itemCategoryService}

	// Shopping List Item Handler
	shoppingListItemHandler := v1.ShoppingListItemHandler{UserRepository: userRepository, ShoppingListRepository: shoppingListRepository,
		ShoppingListItemRepository: shoppingListItemRepository, ShoppingListItemFactory: shoppingListItemFactory,
		ShoppingListItemImageFactory: shoppingListItemImageFactory}

	// Shopping List Item Image Handler
	shoppingListItemImageHandler := v1.ShoppingListItemImageHandler{UserRepository: userRepository, ShoppingListRepository: shoppingListRepository,
		ShoppingListItemImageService: shoppingListItemImageService, ShoppingListItemImageFactory: shoppingListItemImageFactory,
		ShoppingListItemRepository: shoppingListItemRepository, ShoppingListItemImageRepository: shoppingListItemImageRepository}

	// Deal Cashback Handler
	dealCashbackHandler := v1.DealCashbackHandler{ShoppingListRepository: shoppingListRepository, DealCashbackService: dealCashbackService}

	// Deal Handler
	dealHandler := v1.DealHandler{DealService: dealService, DealTransformer: dealTransformer, ItemCategoryService: itemCategoryService,
		DealCashbackService: dealCashbackService, UserRepository: userRepository, ItemSubCategoryRepository: itemSubCategoryRepository}

	// Grocer Handler
	grocerHandler := v1.GrocerHandler{GrocerRepository: grocerRepository, GrocerTransformer: grocerTransformer}

	// Event Handler
	eventHandler := v1.EventHandler{EventService: eventService}

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

		// Grocer Routes
		version1.GET("/grocers", grocerHandler.Index)

		// Deal Handler
		version1.GET("deals", dealHandler.ViewAllForGuestUser)

		// Protected Routes
		version1.Use(middlewares.Auth())
		{
			// User Routes
			version1.PATCH("/users/:guid", userHandler.Update)
			version1.GET("/users/:guid", userHandler.View)

			// Device Routes
			version1.DELETE("/devices/:uuid", deviceHandler.Delete)

			// Authentication Routes
			version1.GET("/auth/refresh", authHandler.Refresh)
			version1.GET("/auth/logout", authHandler.Logout)

			// Shopping List Routes
			version1.GET("users/:guid/shopping_lists", shoppingListHandler.View)
			version1.POST("users/:guid/shopping_lists", shoppingListHandler.Create)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Update)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid", shoppingListHandler.Delete)

			// Shopping List Item Routes
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.View)
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.ViewAll)
			version1.POST("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.Create)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Update)
			version1.PATCH("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.UpdateAll)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid", shoppingListItemHandler.Delete)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items", shoppingListItemHandler.DeleteAll)

			// Shopping List Item Image Routes
			version1.GET("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guid", shoppingListItemImageHandler.View)
			version1.POST("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images", shoppingListItemImageHandler.Create)
			version1.DELETE("users/:guid/shopping_lists/:shopping_list_guid/items/:item_guid/images/:image_guids", shoppingListItemImageHandler.Delete)

			// Deal Cashback Handler
			version1.POST("users/:guid/deal_cashbacks", dealCashbackHandler.Create)

			// Deal Handler
			version1.GET("users/:guid/deals", dealHandler.ViewAllForRegisteredUser)
			version1.GET("deals/:deal_guid", dealHandler.View)
			version1.GET("users/:guid/deals/categories", dealHandler.ViewAndGroupByCategory)
			version1.GET("users/:guid/deals/categories/:category_guid", dealHandler.ViewByCategory)
			version1.GET("users/:guid/deals/categories/:category_guid/subcategories", dealHandler.ViewByCategoryAndGroupBySubCategory)
			version1.GET("users/:guid/deals/subcategories/:subcategory_guid", dealHandler.ViewBySubCategory)

			// Feature Deal (In Carousel) Handler
			version1.GET("users/:guid/featured_deals", eventHandler.ViewAll)
		}
	}

	return router
}
