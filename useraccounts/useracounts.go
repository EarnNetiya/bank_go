package useraccounts

import (
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"goproject-bank/transactions"
	"strconv"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount{
	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}

	database.DB.Where("id = ?", id).First(&account)
	account.Balance = uint(amount)
	database.DB.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)
	return responseAcc
}

func getAccount(id uint) *interfaces.Account{
	account := &interfaces.Account{}
	if database.DB.Where("id = ?", id).First(&account).RecordNotFound() {
		return nil
	}
	return account
}

func Transactions(userId uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
    userIdString := strconv.Itoa(int(userId))
    isValid := helpers.ValidateToken(userIdString, jwt)

    if !isValid {
        return map[string]interface{}{"message": "not valid values"}
    }

    fromAccount := getAccount(from)
    toAccount := getAccount(to)

    if fromAccount == nil || toAccount == nil {
        return map[string]interface{}{"message": "Account not found"}
    }

    if fromAccount.UserID != userId {
        return map[string]interface{}{"message": "You are not owner of the account"}
    }

    if int(fromAccount.Balance) < amount {
        return map[string]interface{}{"message": "Not enough balance"}
    }

    updatedFrom := updateAccount(from, int(fromAccount.Balance)-amount)
    updatedTo := updateAccount(to, int(toAccount.Balance)+amount)

    transactions.CreateTransaction(from, to, amount)

    return map[string]interface{}{
        "message": "all is fine",
        "data":    map[string]interface{}{"from": updatedFrom, "to": updatedTo},
    }
}
