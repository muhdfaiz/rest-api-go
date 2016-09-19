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
func (database *Database) setConfigs() {
	configs := &Configs{}
	dbHost := configs.Get("database.yaml", "db_host", "localhost")
	dbPort := configs.Get("database.yaml", "db_port", "3306")
	dbUsername := configs.Get("database.yaml", "db_username", "root")
	dbPassword := configs.Get("database.yaml", "db_password", "")
	dbName := configs.Get("database.yaml", "db_name", "shoppermate-api")

	database.host = dbHost
	database.port = dbPort
	database.username = dbUsername
	database.password = dbPassword
	database.name = dbName
}

// Connect function will connect with the database
func (database *Database) Connect() *gorm.DB {
	database.setConfigs()

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
