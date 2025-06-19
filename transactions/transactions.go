package transactions

import (
	// "goproject-bank/helpers"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"goproject-bank/blockchain"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

func CreateTransactionByAccountNumbers(fromAccountNumber string, toAccountNumber string, amount int, userID uint) map[string]interface{} {
    db := database.DB
    log.Println("Creating transaction:", fromAccountNumber, "->", toAccountNumber, "Amount:", amount)

    // ตรวจสอบบัญชีต้นทาง
    log.Printf("Checking account: AccountNumber=%s, UserID=%d", fromAccountNumber, userID)
    
    var fromAccount interfaces.Account
    if err := db.Where("account_number = ? AND user_id = ?", fromAccountNumber, userID).First(&fromAccount).Error; err != nil {
        log.Println("From account not found:", err)
        return map[string]interface{}{"message": "From account not found"}
    }
    log.Printf("Found account: ID=%d, Balance=%d", fromAccount.ID, fromAccount.Balance)

    // ตรวจสอบบัญชีปลายทาง
    var toAccount interfaces.Account
    if err := db.Where("account_number = ?", toAccountNumber).First(&toAccount).Error; err != nil {
        log.Println("To account not found:", err)
        return map[string]interface{}{"message": "To account not found"}
    }

    if fromAccount.Balance < uint(amount) {
        log.Println("Insufficient balance in account", fromAccountNumber)
        return map[string]interface{}{"message": "Insufficient balance"}
    }

    // คำนวณแฮชสำหรับธุรกรรม
    transactionData := interfaces.BlockchainTransaction{
        SenderAccount:   fromAccountNumber,
        ReceiverAccount: toAccountNumber,
        Amount:          float64(amount),
        Timestamp:       time.Now().String(),
    }
    dataBytes, _ := json.Marshal(transactionData)
    hash := sha256.Sum256(dataBytes)
    transactionHash := fmt.Sprintf("%x", hash)

    // บันทึกธุรกรรมในฐานข้อมูล
    err := db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-uint(amount)).Error; err != nil {
            return err
        }
        if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+uint(amount)).Error; err != nil {
            return err
        }

        newTransaction := interfaces.Transactions{
            FromAccountNumber: fromAccountNumber,
            ToAccountNumber:   toAccountNumber,
            Amount:            amount,
            Timestamp:         time.Now(),
            Hash:              transactionHash,
        }
        if err := tx.Create(&newTransaction).Error; err != nil {
            return err
        }

        // add to blockchain
        blockchain.Chain.AddTransaction(fromAccountNumber, toAccountNumber, float64(amount))

        return nil
    })

    if err != nil {
        log.Println("Transaction failed:", err)
        return map[string]interface{}{"message": "Transaction failed", "error": err.Error()}
    }

    log.Println("Transaction successful:", fromAccountNumber, "->", toAccountNumber, "Amount:", amount)
    return map[string]interface{}{"message": "Transaction successful"}
}


func GetTransactionsByAccount(accountNumber string) []interfaces.Transactions {
	transactions := []interfaces.Transactions{}
	database.DB.Where("from_account_number = ? OR to_account_number = ?", accountNumber, accountNumber).Find(&transactions)
	return transactions
}


func GetMyTransactions(id string, jwt string) map[string]interface{} {
    userID, err := helpers.ExtractUserID(jwt)
    if err != nil || fmt.Sprintf("%d", userID) != id {
        log.Println("Invalid token or user ID mismatch:", err)
        return map[string]interface{}{"message": "not valid values"}
    }
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
