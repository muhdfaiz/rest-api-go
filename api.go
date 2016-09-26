package main

import (
	"bitbucket.org/cliqers/shoppermate-api/application"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Router
	router := gin.New()
	router.Use(gin.Recovery())

	Database := &systems.Database{}
	db := Database.Connect()

	application.Bootstrap(application.InitializeObjectAndSetRoutes(router, db)).Run(":8080")
}
