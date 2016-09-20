package application

import (
	"bitbucket.org/shoppermate-api/application/v1"
	"bitbucket.org/shoppermate-api/middlewares"
	"github.com/gin-gonic/gin"
)

// SetRoutes will set all routes across the API
func SetRoutes(router *gin.Engine) {
	smsHandler := v1.SmsHandler{}
	userHandler := v1.UserHandler{}
	deviceHandler := v1.DeviceHandler{}
	authHandler := v1.AuthHandler{}

	router.Use(middlewares.Loader())

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
		}

	}
}
