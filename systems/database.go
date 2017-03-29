package systems

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// Database will handle database connection initiliazation.
type Database struct {
	host     string
	port     string
	username string
	password string
	name     string
}

// SetConfigs will retrieve database config from .env and set all of the configs.
// If environment parameter equal to `test`, API will used test database for the configs.
// Otherwise API will use production database.
func (database *Database) setConfigs(environment string) {
	switch environment {
	case "test":
		database.host = os.Getenv("TEST_DB_HOST")
		database.port = os.Getenv("TEST_DB_PORT")
		database.username = os.Getenv("TEST_DB_USERNAME")
		database.password = os.Getenv("TEST_DB_PASSWORD")
		database.name = os.Getenv("TEST_DB_NAME")
	default:
		database.host = os.Getenv("DB_HOST")
		database.port = os.Getenv("DB_PORT")
		database.username = os.Getenv("DB_USERNAME")
		database.password = os.Getenv("DB_PASSWORD")
		database.name = os.Getenv("DB_NAME")
	}
}

// Connect function will initiate new connection with the database.
func (database *Database) Connect(environment string) *gorm.DB {
	database.setConfigs(environment)

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		database.username, database.password, database.host, database.port, database.name))

	if err != nil {
		panic(err)
	}

	err = db.DB().Ping()

	if err != nil {
		panic(err)
	}

	db = database.SetLogger(db)
	db = database.SetConnection(db, 0, 100)

	return db
}

// SetConnection function used to customize database connection settings.
// For example want to increate max idle connection and max open connection allowed.
func (database *Database) SetConnection(db *gorm.DB, maxIdleConnection, maxOpenConnection int) *gorm.DB {
	db.DB().SetMaxIdleConns(maxIdleConnection)
	db.DB().SetMaxOpenConns(maxOpenConnection)

	return db
}

// SetLogger function used to print out any mysql query to stdout.
func (database *Database) SetLogger(db *gorm.DB) *gorm.DB {
	if os.Getenv("DEBUG_DATABASE") == "true" {
		db.SetLogger(gorm.Logger{revel.TRACE})
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
		db.LogMode(true)
	}

	return db
}
