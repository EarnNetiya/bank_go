package users
import (
    "math/rand"
    "time"
)

func generateRandomAccountNumber() string {
    rand.Seed(time.Now().UnixNano())
    var digits = "0123456789"
    accountNum := make([]byte, 10)
    for i := range accountNum {
        accountNum[i] = digits[rand.Intn(len(digits))]
    }
    return string(accountNum)
}
