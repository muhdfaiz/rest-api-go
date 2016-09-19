package application

import (
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/shoppermate-api/systems"

	"github.com/gin-gonic/gin"
)

// Bootstrap func will initialize application with default configurations
func Bootstrap() *gin.Engine {
	fileErrorLog, err := os.OpenFile(fmt.Sprintf("%s/shoppermate-api_error.log", os.Getenv("GOPATH")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	defer fileErrorLog.Close()

	gin.DefaultWriter = io.MultiWriter(fileErrorLog, os.Stdout)
	gin.SetMode(gin.ReleaseMode)

	config := systems.Configs{}

	if config.Get("app.yaml", "debug", "") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Router
	router := gin.New()

	router.Use(gin.LoggerWithWriter(io.MultiWriter(fileErrorLog, os.Stdout)))
	router.Use(gin.Recovery())

	// Load all routes
	SetRoutes(router)

	return router
}
