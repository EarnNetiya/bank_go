package main

import (
	// "goproject-bank/api"
	"goproject-bank/migrations"
	// "goproject-bank/migrations"
)

func main() {
	migrations.MigrateTransactions()
	// api.StartApi()
}