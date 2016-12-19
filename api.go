package main

import (
	"log"
	"reflect"
	"sync"

	"bitbucket.org/cliqers/shoppermate-api/application"
	"bitbucket.org/cliqers/shoppermate-api/systems"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	validator "gopkg.in/go-playground/validator.v8"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	binding.Validator = new(defaultValidator)

	Database := &systems.Database{}
	DB := Database.Connect("production")

	// Initialize Router
	router := gin.New()
	router.Use(gin.Recovery())
	application.Bootstrap(application.InitializeObjectAndSetRoutes(router, DB)).Run(":8080")
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}
	return nil
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		config := &validator.Config{TagName: "binding", FieldNameTag: "json"}
		v.validate = validator.New(config)
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()
	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
