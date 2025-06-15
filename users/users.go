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
	// Sign token
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("JWT_SECRET")) // แก้ไขชื่อ secret ให้ถูกต้อง
	if err != nil {
		return "", err
	}
	return token, nil
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount, withToken bool) (map[string]interface{}, error) {
	// Setup response
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

var secretKey = []byte("supersecretkey") // ต้องตรงกับ helpers.secretKey

func Login(username, password string) map[string]interface{} {
	var user interfaces.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return map[string]interface{}{"message": "User not found"}
		}
		return map[string]interface{}{"message": "Database error", "error": err.Error()}
	}

	// ตรวจสอบรหัสผ่าน
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return map[string]interface{}{"message": "Invalid password"}
	}

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return map[string]interface{}{"message": "Failed to generate token", "error": err.Error()}
	}

	return map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
		"user_id": user.ID,
	}
}

func Register(username string, email string, pass string, initialAmount int) map[string]interface{} {
	// Validate input
	if !helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		},
	) {
		return map[string]interface{}{"message": "not valid values"}
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
		Balance:       uint(initialAmount),
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

func GetUser(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)

	if !isValid {
		return map[string]interface{}{"message": "Invalid token"}
	}

	user := interfaces.User{}
	result := database.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return map[string]interface{}{"message": "User not found"}
	}

	var accounts []interfaces.ResponseAccount
	database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	response, err := prepareResponse(&user, accounts, false)
	if err != nil {
		return map[string]interface{}{"message": "Failed to prepare response"}
	}
	return response
}
