package transactions

import (
	// "goproject-bank/helpers"
	"goproject-bank/blockchain"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"time"

	"github.com/jinzhu/gorm"
)

func CreateTransactionByAccountNumbers(fromAccountNumber string, toAccountNumber string, amount int) map[string]interface{} {

	var fromAccount interfaces.Account
	if err := database.DB.Where("account_number = ?", fromAccountNumber).First(&fromAccount).Error; err != nil {
		return map[string]interface{}{"message": "From account not found"}
	}

	var toAccount interfaces.Account
	if err := database.DB.Where("account_number = ?", toAccountNumber).First(&toAccount).Error; err != nil {
		return map[string]interface{}{"message": "To account not found"}
	}

	// Check balance
	if fromAccount.Balance < uint(amount) {
		return map[string]interface{}{"message": "Insufficient balance"}
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-uint(amount)).Error; err != nil {
			return err
		}

		// Add to toAccount
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+uint(amount)).Error; err != nil {
			return err
		}

		transaction := &interfaces.Transactions{From: fromAccount.ID, To: toAccount.ID, Amount: amount}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		blockData := interfaces.BlockchainTransaction{
			SenderAccount:   fromAccountNumber,
			ReceiverAccount: toAccountNumber,
			Amount:          float64(amount),
			Timestamp:       time.Now().Format(time.RFC3339),
		}
		blockchain.Chain.AddBlock(blockData)

		return nil
	})

	if err != nil {
		return map[string]interface{}{"message": "Transaction failed", "error": err.Error()}
	}

	return map[string]interface{}{"message": "Transaction successful"}
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