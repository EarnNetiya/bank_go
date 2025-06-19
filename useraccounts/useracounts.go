package useraccounts

import (
	"goproject-bank/blockchain"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"strconv"
	"time"
)

func CreateAccount(userId uint, name string, Balance int) (interfaces.ResponseAccount, error) {
	account := interfaces.Account{
		UserID:  userId,
		Name:    name,
		Balance: uint(Balance),
	}
	if err := database.DB.Create(&account).Error; err != nil {
		return interfaces.ResponseAccount{}, err
	}

	return interfaces.ResponseAccount{
		ID:      account.ID,
		Name:    account.Name,
		Balance: int(account.Balance),
	}, nil
}

func Transactions(userId uint, fromAccNumber string, toAccNumber string, amount int, jwt string) map[string]interface{} {
	userIdString := strconv.Itoa(int(userId))
	isValid := helpers.ValidateUserToken(userIdString, jwt)

	if !isValid {
		return map[string]interface{}{"message": "Invalid token"}
	}

	fromAccount, errFrom := GetAccountByNumber(fromAccNumber)
	toAccount, errTo := GetAccountByNumber(toAccNumber)

	if errFrom != nil || errTo != nil {
		return map[string]interface{}{"message": "Account not found"}
	}

	if fromAccount.UserID != userId {
		return map[string]interface{}{"message": "You are not owner of the account"}
	}

	if int(fromAccount.Balance) < amount {
		return map[string]interface{}{"message": "Not enough balance"}
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&interfaces.Account{}).Where("id = ?", fromAccount.ID).
		Update("balance", fromAccount.Balance-uint(amount)).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to update from account"}
	}

	if err := tx.Model(&interfaces.Account{}).Where("id = ?", toAccount.ID).
		Update("balance", toAccount.Balance+uint(amount)).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to update to account"}
	}

	transaction := interfaces.Transactions{
		FromAccountNumber:   fromAccount.AccountNumber,
		ToAccountNumber:     toAccount.AccountNumber,
		Amount: amount,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to create transaction"}
	}

	blockData := interfaces.BlockchainTransaction{
		SenderAccount:   fromAccount.AccountNumber,
		ReceiverAccount: toAccount.AccountNumber,
		Amount:          float64(amount),
		Timestamp:       time.Now().Format(time.RFC3339),
	}
	blockchain.Chain.AddBlock(blockData)

	if err := tx.Commit().Error; err != nil {
		return map[string]interface{}{"message": "Transaction commit failed"}
	}

	respFrom := interfaces.ResponseAccount{
		ID:      fromAccount.ID,
		Name:    fromAccount.Name,
		Balance: int(fromAccount.Balance - uint(amount)),
	}
	respTo := interfaces.ResponseAccount{
		ID:      toAccount.ID,
		Name:    toAccount.Name,
		Balance: int(toAccount.Balance + uint(amount)),
	}

	return map[string]interface{}{
		"message": "Transaction successful",
		"data": map[string]interface{}{
			"from": respFrom,
			"to":   respTo,
		},
	}
}


func GetAccountByNumber(accountNumber string) (*interfaces.Account, error) {
	account := &interfaces.Account{}
	err := database.DB.Where("account_number = ?", accountNumber).First(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}
