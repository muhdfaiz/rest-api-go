package application

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Bootstrap func will initialize application with default configurations
func Bootstrap(router *gin.Engine) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	return router
}
