package interfaces

import ("github.com/jinzhu/gorm")

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
	AccountNum string
}

type ResponseTransaction struct {
	ID uint
	From uint
	To uint
	Amount int
}

type ResponseAccount struct {
	ID uint
	Name string
	Balance int
}

type ResponseUser struct {
	ID uint
	Username string
	Email    string
	Accounts []ResponseAccount
}

type Validation struct {
	Value string
	Valid string
}

type ErrResponse struct {
	Message string
}

type Transactions struct {
	gorm.Model
	From uint
	To uint
	Amount int
}

type AdminOnly struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type ResponseAdmin struct {
	ID       uint
	Username string
	Email    string
}

type BlockchainTransaction struct {
	SenderAccount   string
	ReceiverAccount string
	Amount          float64
	Timestamp       string
}