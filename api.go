package main

import (
	"log"
	"os"

	"bitbucket.org/cliqers/shoppermate-api/application/v1"
	"bitbucket.org/cliqers/shoppermate-api/application/v11"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.DebugMode)

	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	binding.Validator = new(systems.DefaultValidator)

	Database := &systems.Database{}
	DB := Database.Connect("production")

	// Initialize Router
	router := gin.Default()
	routerForSSL := gin.Default()

	if os.Getenv("ENABLE_HTTPS") == "true" {
		routerForSSL = v1.InitializeObjectAndSetRoutesV1(routerForSSL, DB)
		routerForSSL = v11.InitializeObjectAndSetRoutesV1_1(routerForSSL, DB)
		go routerForSSL.RunTLS(":8081", os.Getenv("FULLCHAIN_KEY"), os.Getenv("PRIVATE_KEY"))
	}

	router = v1.InitializeObjectAndSetRoutesV1(router, DB)
	router = v11.InitializeObjectAndSetRoutesV1_1(router, DB)
	router.Run(":8080")
}
