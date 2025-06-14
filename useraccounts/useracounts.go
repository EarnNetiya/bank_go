package useraccounts

import (
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"strconv"
)

func CreateAccount(userId uint, name string, initialAmount int) (interfaces.ResponseAccount, error) {
	account := interfaces.Account{
		UserID:  userId,
		Name:    name,
		Balance: uint(initialAmount),
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

func updateAccount(id uint, amount int) (interfaces.ResponseAccount, error) {
	account := interfaces.Account{}
	if err := database.DB.Where("id = ?", id).First(&account).Error; err != nil {
		return interfaces.ResponseAccount{}, err
	}
	account.Balance = uint(amount)
	if err := database.DB.Save(&account).Error; err != nil {
		return interfaces.ResponseAccount{}, err
	}

	return interfaces.ResponseAccount{
		ID:      account.ID,
		Name:    account.Name,
		Balance: int(account.Balance),
	}, nil
}

func getAccount(id uint) (*interfaces.Account, error) {
	account := &interfaces.Account{}
	err := database.DB.Where("id = ?", id).First(account).Error
	if err != nil {
		return nil, err
	}
	return account, nil
}

func Transactions(userId uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIdString := strconv.Itoa(int(userId))
	isValid := helpers.ValidateToken(userIdString, jwt)

	if !isValid {
		return map[string]interface{}{"message": "Invalid token"}
	}

	fromAccount, errFrom := getAccount(from)
	toAccount, errTo := getAccount(to)

	if errFrom != nil || errTo != nil {
		return map[string]interface{}{"message": "Account not found"}
	}

	if fromAccount.UserID != userId {
		return map[string]interface{}{"message": "You are not owner of the account"}
	}

	if int(fromAccount.Balance) < amount {
		return map[string]interface{}{"message": "Not enough balance"}
	}

	// เริ่ม transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&interfaces.Account{}).Where("id = ?", from).Update("balance", fromAccount.Balance-uint(amount)).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to update from account"}
	}

	if err := tx.Model(&interfaces.Account{}).Where("id = ?", to).Update("balance", toAccount.Balance+uint(amount)).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to update to account"}
	}

	transaction := interfaces.Transactions{
		From:   from,
		To:     to,
		Amount: amount,
	}
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return map[string]interface{}{"message": "Failed to create transaction"}
	}

	if err := tx.Commit().Error; err != nil {
		return map[string]interface{}{"message": "Transaction commit failed"}
	}

	updatedFrom, err := getAccount(from)
	if err != nil {
		return map[string]interface{}{"message": "Failed to get updated from account"}
	}

	updatedTo, err := getAccount(to)
	if err != nil {
		return map[string]interface{}{"message": "Failed to get updated to account"}
	}

	respFrom := interfaces.ResponseAccount{
		ID:      updatedFrom.ID,
		Name:    updatedFrom.Name,
		Balance: int(updatedFrom.Balance),
	}
	respTo := interfaces.ResponseAccount{
		ID:      updatedTo.ID,
		Name:    updatedTo.Name,
		Balance: int(updatedTo.Balance),
	}

	return map[string]interface{}{
		"message": "Transaction successful",
		"data": map[string]interface{}{
			"from": respFrom,
			"to":   respTo,
		},
	}
}
