package transaction

import (
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	db := helpers.ConnectDB()
	transaction := &interfaces.Transactions{From: From, To: To, Amount: Amount}
	db.Create(transaction)

	defer db.Close()
}