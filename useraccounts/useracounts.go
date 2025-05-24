package useraccounts

import (
	"fmt"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
)

func updateAccount(id uint, amount int) {
	db := helpers.ConnectDB()
	db.Model(&interfaces.Account{}).Where("id = ?", id).Update("balance", amount)
	defer db.Close()
}

func getAccount(id uint) *interfaces.Account{
	db := helpers.ConnectDB()
	account := &interfaces.Account{}
	if db.Where("id = ?", id).First(&account).RecordNotFound() {
		return nil
	}
	defer db.Close()
	return account
}

func Transactions(userId uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIdString := fmt.Strint(userId)
	isValid := helpers.ValidateToken(userIdString, jwt)

	if isValid {

	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}