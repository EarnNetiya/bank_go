package main

import (
	"goproject-bank/api"
	"goproject-bank/database"
	// "goproject-bank/migrations"
)

func main() {
	// migrations.MigrateTransactions()
	database.InitDatabase()
	api.StartApi()
}