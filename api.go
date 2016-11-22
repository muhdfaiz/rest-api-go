package main

import (
	"bitbucket.org/cliqers/shoppermate-api/application"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

func main() {
	Database := &systems.Database{}
	DB := Database.Connect("production")

	// Initialize Router
	router := gin.New()
	router.Use(gin.Recovery())
	application.Bootstrap(application.InitializeObjectAndSetRoutes(router, DB)).Run(":8080")
}
