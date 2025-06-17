package admin

import (
	"fmt"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CreateAccount(account *interfaces.Account) (bool, string) {
	var existingAccount interfaces.Account
	if !database.DB.Where("account_number = ?", account.AccountNumber).First(&existingAccount).RecordNotFound() {
		return false, "Account number already exists"
	}

	if err := database.DB.Create(account).Error; err != nil {
		return false, "Failed to create account"
	}

	return true, ""
}

func CreateAdmin(admin *interfaces.AdminOnly) bool {
	if err := database.DB.Create(admin).Error; err != nil {
		return false
	}
	return true
}

func prepareAdminToken(admin *interfaces.AdminOnly) string {
	tokenContent := jwt.MapClaims{
		"admin_id":          admin.ID,
		"admin":             true,
		"verification_code": helpers.AdminVerificationCode,
		"exp":               time.Now().Add(time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenContent)
	token, err := jwtToken.SignedString(helpers.JwtSecret)
	helpers.HandleErr(err)
	return token
}


func prepareResponse(admin *interfaces.AdminOnly, users []interfaces.ResponseUser, withToken bool) map[string]interface{} {
	responseUser := &interfaces.ResponseUser{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
	}

	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		token := prepareAdminToken(admin)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response
}

func Login(username, pass string) map[string]interface{} {
	
	valid := helpers.ValidationAdmin([]interfaces.Validation{
		{Value: username, Valid: "username"},
		{Value: pass, Valid: "password"},
	})
	if !valid {
		return map[string]interface{}{"message": "not valid values"}
	}

	admin := &interfaces.AdminOnly{}
	if database.DB.Where("username = ?", username).First(&admin).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"}
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(pass))
	if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
		return map[string]interface{}{"message": "Wrong password"}
	}

	claims := jwt.MapClaims{
        "admin": true,
        "verification_code": helpers.AdminVerificationCode, // Assuming this exists
        "exp":   time.Now().Add(time.Hour * 24).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(helpers.JwtSecret)
    if err != nil {
        return map[string]interface{}{"message": "Failed to generate token"}
    }
    return map[string]interface{}{"message": "Login successful", "token": tokenString}
}

func Register(username, email, password string) map[string]interface{} {
	valid := helpers.ValidationAdmin([]interfaces.Validation{
		{Value: username, Valid: "username"},
		{Value: email, Valid: "email"},
		{Value: password, Valid: "password"},
	})

	if !valid {
		return map[string]interface{}{"message": "Invalid registration data"}
	}

	hashedPassword := helpers.HashAndSalt([]byte(password))
	fmt.Println("Registering admin:", username, email, hashedPassword)

	admin := interfaces.AdminOnly{
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	success := CreateAdmin(&admin)
	if !success {
		fmt.Println("Unable to register admin in DB")
		return map[string]interface{}{"message": "Unable to register admin"}
	}

	token := prepareAdminToken(&admin)
	return map[string]interface{}{
		"message": "Admin registered successfully",
		"token":   token,
	}
}

func GetAllUser(id, auth string) map[string]interface{} {
	if !helpers.ValidateAdminToken(auth) {
		return map[string]interface{}{"message": "Invalid token"}
	}
	idUint, _ := strconv.ParseUint(id, 10, 64)

	admin := interfaces.AdminOnly{}
	if database.DB.Where("id = ?", idUint).First(&admin).RecordNotFound() {
		return map[string]interface{}{"message": "Admin not found"}
	}

	var users []interfaces.ResponseUser
	database.DB.Table("users").Select("id, username, email").Scan(&users)
	return prepareResponse(&admin, users, false)
}

func GetUser(id, auth string) map[string]interface{} {
	log.Println("GetUser: id =", id)
	log.Println("Token =", auth)

	if !helpers.ValidateAdminToken(auth) {
		log.Println("Token failed validation")
		return map[string]interface{}{"message": "Unauthorized", "status": 401}
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Invalid user ID format:", err)
		return map[string]interface{}{"message": "Invalid user ID", "status": 400}
	}

	user := interfaces.User{}
	if database.DB.Where("id = ?", idUint).First(&user).RecordNotFound() {
		log.Println("User not found")
		return map[string]interface{}{"message": "User not found", "status": 404}
	}

	var accounts []interfaces.ResponseAccount
	database.DB.Table("accounts").Select("id, account_number, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	log.Println("User retrieved successfully")
	return map[string]interface{}{
		"message": "User retrieved successfully",
		"status":  200,
		"data": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"accounts": accounts,
		},
	}
}

func GetAccount(id, auth string) map[string]interface{} {
	log.Println("GetAccount: id =", id)
	log.Println("Token =", auth)

	if !helpers.ValidateAdminToken(auth) {
		log.Println("Token failed validation")
		return map[string]interface{}{"message": "Unauthorized", "status": 401}
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Invalid account ID format:", err)
		return map[string]interface{}{"message": "Invalid account ID", "status": 400}
	}

	account := interfaces.Account{}
	if database.DB.Where("id = ?", idUint).First(&account).RecordNotFound() {
		log.Println("Account not found")
		return map[string]interface{}{"message": "Account not found", "status": 404}
	}

	log.Println("Account retrieved successfully")
	return map[string]interface{}{
		"message": "Account retrieved successfully",
		"status":  200,
		"data": map[string]interface{}{
			"id":             account.ID,
			"account_number": account.AccountNumber,
			"name":           account.Name,
			"balance":        account.Balance,
			"user_id":        account.UserID,
		},
	}
}


func DeleteUser(id, auth string) map[string]interface{} {
	log.Println("DeleteUser: id =", id)
	log.Println("Token =", auth)

	if !helpers.ValidateAdminToken(auth) {
		log.Println("Token failed validation")
		return map[string]interface{}{"message": "Unauthorized", "status": 401}
	}

	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		log.Println("Invalid user ID format:", err)
		return map[string]interface{}{"message": "Invalid user ID"}
	}

	result := database.DB.Where("id = ?", idUint).Delete(&interfaces.User{})
	if result.RowsAffected == 0 {
		log.Println("No user found or already deleted")
		return map[string]interface{}{"message": "User not found or already deleted"}
	}

	log.Println("User deleted successfully")
	return map[string]interface{}{"message": "User deleted successfully"}
}



func DeleteAccount(id, auth string) map[string]interface{} {
	if !helpers.ValidateAdminToken(auth) {
		return map[string]interface{}{"message": "Invalid token"}
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return map[string]interface{}{"message": "Invalid account ID"}
	}

	if database.DB.Where("id = ?", idUint).Delete(&interfaces.Account{}).RowsAffected == 0 {
		return map[string]interface{}{"message": "Account not found or already deleted"}
	}
	return map[string]interface{}{"message": "Account deleted successfully"}
}
