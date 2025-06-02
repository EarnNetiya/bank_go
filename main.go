package main

import (
	"goproject-bank/api"
	"goproject-bank/database"
	"goproject-bank/migrations"
	// "goproject-bank/migrations"
)

func main() {
	
	database.InitDatabase()
	migrations.Migrate()
	api.StartApi()
}