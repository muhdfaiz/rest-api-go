package application

import (
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
)

// Bootstrap func will initialize application with default configurations
func Bootstrap(router *gin.Engine) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	config := systems.Configs{}

	if config.Get("app.yaml", "debug", "") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	return router
}
