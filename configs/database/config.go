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

type DatabaseTesting struct {
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
			&entities.Group{},
			&entities.Friend{},
			&entities.Transaction{},
			&entities.Payment{},
			&entities.Item{},
			&entities.Activity{},
			&entities.GroupActivity{},
			&entities.PaymentActivity{},
			&entities.TransactionActivity{},
			&entities.ReminderActivity{},
			&entities.PasswordResetTokens{},
			&entities.Expense{},
			&entities.Notification{},
		)

		database.connection = db
	})
}

// Testing Database
func (databaseTesting *DatabaseTesting) lazyInit() {
	databaseTesting.once.Do(func() {
		host := "34.101.52.136"
		port := "5432"
		dbname := "split-rex-db-testing"
		username := "testing_admin"
		password := "kvT9=5@K-Cz*ZIB"

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
			&entities.Group{},
			&entities.Friend{},
			&entities.Transaction{},
			&entities.Payment{},
			&entities.Item{},
			&entities.Activity{},
			&entities.GroupActivity{},
			&entities.PaymentActivity{},
			&entities.TransactionActivity{},
			&entities.ReminderActivity{},
			&entities.PasswordResetTokens{},
			&entities.Expense{},
			&entities.Notification{},
		)

		databaseTesting.connection = db
	})
}

func (database *Database) GetConnection() *gorm.DB {
	database.lazyInit()
	return database.connection
}

func (databaseTesting *DatabaseTesting) GetConnection() *gorm.DB {
	databaseTesting.lazyInit()
	return databaseTesting.connection
}

var DB = &Database{}
var DBTesting = &DatabaseTesting{}
