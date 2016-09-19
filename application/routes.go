package application

import (
	"github.com/gin-gonic/gin"
	"bitbucket.org/shoppermate/application/v1"
	"bitbucket.org/shoppermate/middlewares"
	"bitbucket.org/shoppermate/systems"
)

func Database() gin.HandlerFunc {
	database := &systems.Database{}
	db := database.Connect()
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

// SetRoutes will set all routes across the API
func SetRoutes(router *gin.Engine) {
	smsHandler := v1.SmsHandler{}
	userHandler := v1.UserHandler{}
	deviceHandler := v1.DeviceHandler{}
	authHandler := v1.AuthHandler{}

	router.Use(Database())

	// Simple group: v1
	version1 := router.Group("/v1")
	{

		// Public routes
		// Device Routes
		version1.POST("/devices", deviceHandler.Create)
		version1.PATCH("/devices/:uuid", deviceHandler.Update)

		// User Routes
		// Create User Route
		version1.POST("/users", userHandler.Create)

		// SMS Routes
		// Send SMS Route
		version1.POST("/sms", smsHandler.Send)
		// Verify SMS Verification Code Route
		version1.POST("/sms/verifications", smsHandler.Verify)

		// Authentication Routes
		// Login via phone no Route
		version1.POST("/auth/login/phone", authHandler.LoginViaPhone)
		version1.POST("/auth/login/facebook", authHandler.LoginViaFacebook)

	}

	// Protected Routes
	authentication := router.Group("/v1")

	authentication.Use(middlewares.Auth())
	{
		// User Routes
		// Update User Route
		authentication.PATCH("/users/:guid", userHandler.Update)
		authentication.GET("/users/:guid", userHandler.View)

		// Device Routes
		// Delete Device
		authentication.DELETE("/devices/:uuid", deviceHandler.Delete)

		// Authentication Routes
		// Logout Route
		authentication.GET("/auth/logout", authHandler.Logout)
	}
}
