package interfaces

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
	Accounts []Account `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // One-to-many relationship with Account
}

type Account struct {
    ID            uint      `gorm:"primaryKey"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    UserID        uint      `gorm:"not null"`
    AccountNumber string    `gorm:"unique;not null"`
    Balance       uint      `gorm:"not null"`
    Type          string
    Name          string
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
	AccountNumber string
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
	FromAccountNumber string
    ToAccountNumber   string
    Amount            int
	Timestamp       time.Time `gorm:"not null"`
    Hash            string    `gorm:"not null"`
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
    TransactionID   string  `json:"transaction_id"`
    SenderAccount   string  `json:"sender_account"`
    ReceiverAccount string  `json:"receiver_account"`
    Amount          float64 `json:"amount"`
    Timestamp       string  `json:"timestamp"`
}

type BlockWithHash struct {
    Hash      string                 `json:"hash"`
    PrevHash  string                 `json:"prev_hash"`
    Data      BlockchainTransaction  `json:"data"`
    Timestamp time.Time              `json:"timestamp"`
}