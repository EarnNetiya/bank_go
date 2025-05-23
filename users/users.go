package users

import (
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(username string, pass string) map[string]interface{} {
	// Connect db
	db := helpers.ConnectDB()
	user := interfaces.User{}
	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"}
	}

	// verify password
	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return map[string]interface{}{"message": "Wrong password"}
	}

	// Find account for the user
	account := []interfaces.ResponseAccount{}
	db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&account)


	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: account,
	}

	defer db.Close()

	// Sign token 
	tokenContent := jwt.MapClaims{
		"user_ID": user.ID,
		"expiry": time.Now().Add(time.Minute ^ 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	// Prepare response
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}