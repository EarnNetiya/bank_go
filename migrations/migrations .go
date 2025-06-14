package migrations

import (
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"goproject-bank/users"
	// "math/rand"
	// "strconv"
	// "time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
	AccountNumber string
}


func createAccounts() {

	userList := []interfaces.User{
		// {Username: "Martin", Email: "martin@martin.com"},
		// {Username: "Michael", Email: "michael!michael.com"},
	}

	for _, u := range userList {
		generatedPassword := helpers.HashAndSalt([]byte(u.Username))
		user := interfaces.User{
			Username: u.Username,
			Email:    u.Email,
			Password: generatedPassword,
		}
		if err := database.DB.Create(&user).Error; err != nil {
			panic(err)
		}

		accountNum := users.GenerateRandomAccountNumber()

		// Check for duplicates
		var count int64
		database.DB.Model(&interfaces.Account{}).Where("account_number = ?", accountNum).Count(&count)
		for count > 0 {
			accountNum = users.GenerateRandomAccountNumber()
			database.DB.Model(&interfaces.Account{}).Where("account_number = ?", accountNum).Count(&count)
		}

		account := interfaces.Account{
			Type:          "Daily Account",
			Name:          u.Username + "'s account",
			Balance:       uint(10000 * (len(userList) + 1)),
			UserID:        user.ID,
			AccountNumber: accountNum,
		}
		if err := database.DB.Create(&account).Error; err != nil {
			panic(err)
		}
	}
}

func Migrate() {
	User := &interfaces.User{}
	Account := &interfaces.Account{}
	Transactions := &interfaces.Transactions{}
	Admin := &interfaces.AdminOnly{}
	database.DB.AutoMigrate(&User, &Account, &Transactions, &Admin)

	createAccounts()
}