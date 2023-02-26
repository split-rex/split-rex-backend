package database

import (
	"os"
	"split-rex-backend/entities"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
	once       sync.Once
}

func (database *Database) lazyInit() {
	database.once.Do(func() {
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")
		dbname := os.Getenv("DB_NAME")
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")

		dsn := "host=" + host
		dsn += " user=" + username
		dsn += " password=" + password
		dsn += " dbname=" + dbname
		dsn += " port=" + port
		dsn += " sslmode=disable"

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
		if err != nil {
			panic("Cannot connect database")
		}

		db.AutoMigrate(
			&entities.User{},
			&entities.Friend{},
		)

		database.connection = db
	})
}

func (database *Database) GetConnection() *gorm.DB {
	database.lazyInit()
	return database.connection
}

var DB = &Database{}
