package core

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database is model
type Database struct {
	GormDB *gorm.DB
}

// DatabaseCreator is databse create
type DatabaseCreator struct {
	config DB
}

// CreateDialector is create gorm dialector
func (creator *DatabaseCreator) CreateDialector() gorm.Dialector {
	var dialector gorm.Dialector
	dbtype := strings.ToLower(creator.config.Type)

	switch dbtype {
	case "mysql":
		var dsn string
		if creator.config.Protocol == "unix" {
			dsn = fmt.Sprintf(
				"%s:%s@%s(%s)/%s?charset=%s&parseTime=true",
				creator.config.User,
				creator.config.Pass,
				creator.config.Protocol,
				creator.config.Host,
				creator.config.Name,
				creator.config.Charset,
			)
		} else {
			dsn = fmt.Sprintf(
				"%s:%s@%s(%s:%d)/%s?charset=%s&parseTime=true",
				creator.config.User,
				creator.config.Pass,
				creator.config.Protocol,
				creator.config.Host,
				creator.config.Port,
				creator.config.Name,
				creator.config.Charset,
			)
		}
		dialector = mysql.Open(dsn)
	default:
		dialector = sqlite.Open(creator.config.Name)
	}

	return dialector
}

// CreateDatabase is create gorm database
func (creator *DatabaseCreator) CreateDatabase(config *gorm.Config) (*gorm.DB, error) {
	return gorm.Open(creator.CreateDialector(), config)

}

// DB is create database/sql.DB
func (db *Database) DB() (*sql.DB, error) {
	return db.GormDB.DB()

}

// NewDatabase is create database
func NewDatabase(config Config) Database {
	creator := DatabaseCreator{
		config: config.DB,
	}

	db, err := creator.CreateDatabase(&gorm.Config{})

	if err != nil {
		log.Fatalf("database connection failed : %v", err)
	}

	if !config.App.Debug {
		db.Logger.LogMode(logger.Silent)
	}

	return Database{
		GormDB: db,
	}
}
