package application

import (
	"bitbucket.org/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

// Bootstrap func will initialize application with default configurations
func Bootstrap() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	config := systems.Configs{}

	if config.Get("app.yaml", "debug", "") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Router
	router := gin.New()
	router.Use(gin.Recovery())

	// Load all routes
	SetRoutes(router)

	return router
}
