package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var JwtSecret []byte
var AdminVerificationCode string

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	JwtSecret = []byte(os.Getenv("JWT_SECRET"))
	AdminVerificationCode = os.Getenv("ADMIN_TOKEN")
}
