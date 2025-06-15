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

func CreateTransactionByAccountNumbers(fromAccountNumber string, toAccountNumber string, amount int, userID uint) map[string]interface{} {
	var fromAccount interfaces.Account
	if err := database.DB.Where("account_number = ?", fromAccountNumber).First(&fromAccount).Error; err != nil {
		return map[string]interface{}{"message": "From account not found"}
	}

	if fromAccount.UserID != userID {
		return map[string]interface{}{"message": "Unauthorized: You do not own this account"}
	}
	
	var toAccount interfaces.Account
	if err := database.DB.Where("account_number = ?", toAccountNumber).First(&toAccount).Error; err != nil {
		return map[string]interface{}{"message": "To account not found"}
	}

	if fromAccount.Balance < uint(amount) {
		return map[string]interface{}{"message": "Insufficient balance"}
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-uint(amount)).Error; err != nil {
			return err
		}
		if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+uint(amount)).Error; err != nil {
			return err
		}

		transaction := &interfaces.Transactions{
			FromAccountNumber: fromAccountNumber,
			ToAccountNumber:   toAccountNumber,
			Amount:            amount,
		}
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




func GetTransactionsByAccount(accountNumber string) []interfaces.Transactions {
	transactions := []interfaces.Transactions{}
	database.DB.Where("from_account_number = ? OR to_account_number = ?", accountNumber, accountNumber).Find(&transactions)
	return transactions
}


func GetMyTransactions(id string, jwt string) map[string]interface{} {
	if helpers.ValidateToken(id, jwt) {
		accounts := []interfaces.Account{}
		database.DB.Where("user_id = ?", id).Find(&accounts)

		var transactions []interfaces.Transactions

		for _, acc := range accounts {
			txs := GetTransactionsByAccount(acc.AccountNumber)
			transactions = append(transactions, txs...)
		}

		return map[string]interface{}{
			"message": "all is fine",
			"data":    transactions,
		}
	}
	return map[string]interface{}{"message": "not valid values"}
}
