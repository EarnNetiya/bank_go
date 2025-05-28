package transactions

import (
	// "goproject-bank/helpers"
	"goproject-bank/blockchain"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"strconv"
	"time"
)

func CreateTransaction(From uint, To uint, Amount int) {
	transaction := &interfaces.Transactions{From: From, To: To, Amount: Amount}
	database.DB.Create(transaction)

	// Add the transaction to the blockchain
	blockData := interfaces.BlockchainTransaction{
		SenderAccount:   strconv.Itoa(int(From)),
		ReceiverAccount: strconv.Itoa(int(To)),
		Amount:          float64(Amount),
		Timestamp:       time.Now().Format(time.RFC3339),
	}

	// Add block in blockchain
	blockchain.Chain.AddBlock(blockData)
}

func GetTransactionsByAccount(id uint) []interfaces.ResponseTransaction {
	transactions := []interfaces.ResponseTransaction{}
	database.DB.Table("transactions").Select("id, transactions.from, transactions.to, amount").Where(interfaces.Transactions{From: id}).Or(interfaces.Transactions{To: id}).Scan(&transactions)
	return transactions
}

func GetMyTransactions(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	if isValid {
		accounts := []interfaces.ResponseAccount{}

		database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", id).Scan(&accounts)

		transactions := []interfaces.ResponseTransaction{}

		for i := 0; i < len(accounts); i++ {
			accTransactions := GetTransactionsByAccount(accounts[i].ID)
			transactions = append(transactions, accTransactions...)
		}
		var response = map[string]interface{}{"message": "all is fine"}
		response["data"] = transactions
		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}

}