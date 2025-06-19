package users

import (

	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"

	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) (string, error) {
    tokenContent := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Minute * 24).Unix(),
    }
    jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
    token, err := jwtToken.SignedString(helpers.JwtSecret)
    if err != nil {
        return "", err
    }
    return token, nil
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) (map[string]interface{}, error) {
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	response := map[string]interface{}{"message": "all is fine"}

	if withToken {
		token, err := prepareToken(user)
		if err != nil {
			return nil, err
		}
		response["jwt"] = token
	}

	response["data"] = responseUser
	return response, nil
}

func Login(email, password string) map[string]interface{} {
	var user interfaces.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return map[string]interface{}{"message": "User not found"}
		}
		return map[string]interface{}{"message": "Database error", "error": err.Error()}
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return map[string]interface{}{"message": "Invalid password"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(helpers.JwtSecret)
	if err != nil {
		return map[string]interface{}{"message": "Failed to generate token", "error": err.Error()}
	}

	return map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"user_id": user.ID,
	}
}

func Register(email string, username string, pass string, Balance int) map[string]interface{} {
	// Validate input
	if !helpers.Validation(
		[]interfaces.Validation{
			{Value: email, Valid: "email"},
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		},
	) {
		return map[string]interface{}{"message": "not valid values"}
	}
	// Check if user already exists
	var existingUser interfaces.User
	err := database.DB.Where("email = ?", email).First(&existingUser).Error
	if err == nil {
		return map[string]interface{}{"message": "Email already exists"}
	} else if err != nil && !gorm.IsRecordNotFoundError(err) {
		return map[string]interface{}{"message": "Database error", "error": err.Error()}
	}

	hashedPass := helpers.HashAndSalt([]byte(pass))

	user := &interfaces.User{
		Username: username,
		Email:    email,
		Password: hashedPass,
	}
	if err := database.DB.Create(user).Error; err != nil {
		return map[string]interface{}{"message": "Unable to create user"}
	}

	accountNum := GenerateRandomAccountNumber()
	account := &interfaces.Account{
		Type:          "Daily Account",
		Name:          username + "'s account",
		Balance:       uint(Balance),
		UserID:        user.ID,
		AccountNumber: accountNum,
	}
	if err := database.DB.Create(account).Error; err != nil {
		return map[string]interface{}{"message": "Unable to create account"}
	}

	respAccount := interfaces.ResponseAccount{
		ID:      account.ID,
		Name:    account.Name,
		Balance: int(account.Balance),
	}
	accounts := []interfaces.ResponseAccount{respAccount}

	response, err := prepareResponse(user, accounts, true)
	if err != nil {
		return map[string]interface{}{"message": "Failed to generate token"}
	}
	return response
}
