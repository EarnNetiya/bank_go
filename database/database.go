package database

import (
	"goproject-bank/helpers"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=db user=postgres password=postgres dbname=gobank port=5432 sslmode=disable"
	database, err := gorm.Open("postgres", dsn)
	helpers.HandleErr(err)

	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)

	DB = database
}