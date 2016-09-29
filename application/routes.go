package application

import (
	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/middlewares"
	"bitbucket.org/cliqers/shoppermate-api/services/facebook"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// InitializeObjectAndSetRoutes will initialize object and set all routes across the API
func InitializeObjectAndSetRoutes(router *gin.Engine, db *gorm.DB) *gin.Engine {
	// User Objects
	userRepository := &v1.UserRepository{DB: db}
	userFactory := &v1.UserFactory{DB: db}
	userService := &v1.UserService{DB: db}

	// Device Objects
	deviceRepository := &v1.DeviceRepository{DB: db}
	deviceFactory := &v1.DeviceFactory{DB: db}

	// Sms Objects
	smsHistoryRepository := &v1.SmsHistoryRepository{DB: db}
	smsService := &v1.SmsService{DB: db}

	// Referral Cashback Objects
	referralCashbackRepository := &v1.ReferralCashbackRepository{DB: db}

	// Facebook Service
	Config := &systems.Configs{}
	facebookService := &facebook.FacebookService{
		AppID:     Config.Get("app.yaml", "facebook_app_id", ""),
		AppSecret: Config.Get("app.yaml", "facebook_app_secret", ""),
	}

	smsHandler := v1.SmsHandler{DB: db, UserRepository: userRepository, UserFactory: userFactory, SmsService: smsService, SmsHistoryRepository: smsHistoryRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}
	userHandler := v1.UserHandler{DB: db, UserRepository: userRepository, UserService: userService, UserFactory: userFactory,
		ReferralCashbackRepository: referralCashbackRepository, SmsService: smsService, FacebookService: facebookService}
	deviceHandler := v1.DeviceHandler{DB: db, UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}
	//shoppingListHandler := v1.ShoppingListHandler{DB: db, UserRepository: userRepository}
	authHandler := v1.AuthHandler{DB: db, UserRepository: userRepository, DeviceRepository: deviceRepository, DeviceFactory: deviceFactory}

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
		version1.Use(middlewares.Auth(db))
		{
			// User Routes
			version1.PATCH("/users/:guid", userHandler.Update)
			version1.GET("/users/:guid", userHandler.View)

			// Shopping List Routes
			//version1.POST("users/:guid/shopping_list", shoppingListHandler.Create)

			// Device Routes
			version1.DELETE("/devices/:uuid", deviceHandler.Delete)

			// Authentication Routes
			version1.GET("/auth/refresh", authHandler.Refresh)
			version1.GET("/auth/logout", authHandler.Logout)
		}

	}

	return router
}
