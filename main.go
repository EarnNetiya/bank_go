package main

import (
	"log"

	"goproject-bank/api"
	"goproject-bank/database"
	"goproject-bank/helpers"       // 👈 สำคัญ!
	"goproject-bank/migrations"
	"goproject-bank/blockchain"
)

func main() {
	helpers.LoadEnv() 

	database.InitDatabase()
	log.Println("Database connected")


	migrations.Migrate()
	log.Println("Migration complete")

	if !blockchain.Chain.VerifyChain() {
		log.Println("Blockchain corrupted. Reinitializing...")
		blockchain.Chain = blockchain.InitBlockChain()
	}
	log.Println("Blockchain initialized")

	api.StartApi()
	log.Println("API started")

	select {} 
}
