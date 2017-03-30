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
	// Load .env file into os environment.
	err := godotenv.Load(os.Getenv("GOPATH") + "src/bitbucket.org/cliqers/shoppermate-api/.env")

	// Return fatal error if cannot load .env file.
	// Fatal error means it will exit this API.
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Always set GIN to run as release mode first.
	// The different between debug mode and release mode is GIN will output any request receive and error in debug mode.
	// Useful for debugging purpose and during development.
	gin.SetMode(gin.ReleaseMode)

	// Set GIN mode to debug mode if `DEBUG` inside .env file equal to true.
	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize validator for request data binding.
	binding.Validator = new(systems.DefaultValidator)

	//Connect to database
	Database := &systems.Database{}
	DB := Database.Connect("production")

	// Initialize Router for SSL and without SSL.
	router := gin.Default()
	routerForSSL := gin.Default()

	// If `ENABLE_HTTPS` setting in .env file equal to true, run this API using HTTPS URI.
	if os.Getenv("ENABLE_HTTPS") == "true" {
		routerForSSL = v1.InitializeObjectAndSetRoutesV1(routerForSSL, DB)
		routerForSSL = v11.InitializeObjectAndSetRoutesV1_1(routerForSSL, DB)
		go routerForSSL.RunTLS(":8081", os.Getenv("FULLCHAIN_KEY"), os.Getenv("PRIVATE_KEY"))
	}

	// Run this API using HTTP URI.
	router = v1.InitializeObjectAndSetRoutesV1(router, DB)
	router = v11.InitializeObjectAndSetRoutesV1_1(router, DB)
	router.Run(":8080")
}
