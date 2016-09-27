package systems

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// Database struct
type Database struct {
	host     string
	port     string
	username string
	password string
	name     string
}

// SetConfigs will retrieve database info from config and initiliaze the config
func (database *Database) setConfigs(environment string) {
	configs := &Configs{}

	switch environment {
	case "test":
		database.host = configs.Get("database.yaml", "db_test_host", "localhost")
		database.port = configs.Get("database.yaml", "db_test_port", "3306")
		database.username = configs.Get("database.yaml", "db_test_username", "root")
		database.password = configs.Get("database.yaml", "db_test_password", "")
		database.name = configs.Get("database.yaml", "db_test_name", "mcliq0621")
	default:
		database.host = configs.Get("database.yaml", "db_host", "localhost")
		database.port = configs.Get("database.yaml", "db_port", "3306")
		database.username = configs.Get("database.yaml", "db_username", "root")
		database.password = configs.Get("database.yaml", "db_password", "")
		database.name = configs.Get("database.yaml", "db_name", "shoppermate")
	}
}

// Connect function will connect with the database
func (database *Database) Connect(environment string) *gorm.DB {
	switch environment {
	case "test":
		database.setConfigs("test")
	default:
		database.setConfigs("production")
	}

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		database.username, database.password, database.host, database.port, database.name))

	if err != nil {
		panic(err)
	}

	// Ping function to checks the database connectivity exist
	err = db.DB().Ping()

	if err != nil {
		panic(err)
	}

	Config := &Configs{}
	if Config.Get("app.yaml", "debug_database", "") == "true" {
		db.SetLogger(gorm.Logger{revel.TRACE})
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
		db.LogMode(true)
	}

	return db
}
