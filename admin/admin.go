package admin

import (
	"fmt"
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
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

	return prepareResponse(admin, nil, true)
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

func DeleteUser(id, auth string) map[string]interface{} {
    if !helpers.ValidateAdminToken(auth) {
        return map[string]interface{}{"message": "Invalid token"}
    }

    fmt.Printf("DeleteUser called with id = %s\n", id) // เช็คค่า id ที่รับมา

    idUint, err := strconv.ParseUint(id, 10, 64)
    if err != nil {
        return map[string]interface{}{"message": "Invalid user ID"}
    }

    if database.DB.Where("id = ?", idUint).Delete(&interfaces.User{}).RowsAffected == 0 {
        return map[string]interface{}{"message": "User not found or already deleted"}
    }
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
