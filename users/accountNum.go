package users

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomAccountNumber() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(1000000000 + rand.Intn(900000000)) // 10 
}