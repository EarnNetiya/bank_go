package useraccounts

import (
	"strconv"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"goproject-bank/transaction"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount{
	db := helpers.ConnectDB()
	account := interfaces.Account{}
	responseAcc := interfaces.ResponseAccount{}

	db.Where("id = ?", id).First(&account)
	account.Balance = uint(amount)
	db.Save(&account)

	responseAcc.ID = account.ID
	responseAcc.Name = account.Name
	responseAcc.Balance = int(account.Balance)
	defer db.Close()
	return responseAcc
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

    transaction.CreateTransaction(from, to, amount)

    return map[string]interface{}{
        "message": "all is fine",
        "data":    map[string]interface{}{"from": updatedFrom, "to": updatedTo},
    }
}
