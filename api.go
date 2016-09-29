package main

import (
	"bitbucket.org/cliqers/shoppermate-api/application"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Router
	router := gin.New()
	router.Use(gin.Recovery())
	application.Bootstrap(application.InitializeObjectAndSetRoutes(router)).Run(":8080")
}
