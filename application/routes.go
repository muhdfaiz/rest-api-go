package application

import (
	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/middlewares"
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

// InitializeObjectAndSetRoutes will initialize object and set all routes across the API
func InitializeObjectAndSetRoutes(router *gin.Engine) *gin.Engine {
	router.Use(middlewares.Loader())

	// User Objects
	userRepository := &v1.UserRepository{}
	userFactory := &v1.UserFactory{}
	userService := &v1.UserService{}

	// Device Objects
	deviceRepository := &v1.DeviceRepository{}
	deviceFactory := &v1.DeviceFactory{}

	// Sms Objects
	smsHistoryRepository := &v1.SmsHistoryRepository{}
	smsService := &v1.SmsService{}

	// Referral Cashback Objects
	referralCashbackRepository := &v1.ReferralCashbackRepository{}

	// Facebook Service
	Config := &systems.Configs{}
	facebookService := &facebook.FacebookService{
		AppID:     Config.Get("app.yaml", "facebook_app_id", ""),
		AppSecret: Config.Get("app.yaml", "facebook_app_secret", ""),
	}

	smsHandler := v1.SmsHandler{UserRepository: userRepository, UserFactory: userFactory, SmsService: smsService, SmsHistoryRepository: smsHistoryRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}
	userHandler := v1.UserHandler{UserRepository: userRepository, UserService: userService, UserFactory: userFactory, DeviceFactory: deviceFactory,
		ReferralCashbackRepository: referralCashbackRepository, SmsService: smsService, FacebookService: facebookService}
	deviceHandler := v1.DeviceHandler{UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}
	//shoppingListHandler := v1.ShoppingListHandler{UserRepository: userRepository}
	authHandler := v1.AuthHandler{UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory, SmsService: smsService}
	//occasionHandler := v1.OccasionHandler{DB: db}

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
			//version1.POST("users/:guid/shopping_list", shoppingListHandler.Create)

			// Occasion Routes
			//version1.GET("/occasions", occasionHandler.Index)
		}

	}

	return router
}
