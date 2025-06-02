package migrations

import (
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"

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
}

func createAccounts() {

	users := &[2]interfaces.User{
		// {Username: "Martin", Email: "martin@martin.com"},
		// {Username: "Michael", Email: "michael!michael.com"},
	}

	for i := 0; i < len(users); i++ {
		generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
		user := User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
		database.DB.Create(&user)

		accout := Account{Type: "Daily Account", Name: string(users[i].Username + "'s" + " accout"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
		database.DB.Create(&accout)
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