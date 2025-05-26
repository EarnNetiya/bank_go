package admin

import (
	"goproject-bank/database"
	"goproject-bank/helpers"
	"goproject-bank/interfaces"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func prepareAdminToken(admin *interfaces.AdminOnly) string {
	// Sign token 
	tokenContent := jwt.MapClaims{
		"admin_id": admin.ID,
		"admin": true,
		"expiry": time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

func prepareResponse(admin *interfaces.AdminOnly, users []interfaces.ResponseUser, withToken bool) map[string]interface{} {
	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
	}

	// Prepare response
	var response = map[string]interface{}{"message": "all is fine"}
	if withToken {
		token := prepareAdminToken(admin)
		response["jwt"] = token
	}
	response["data"] = responseUser

	return response

}

func Login(username string, pass string) map[string]interface{} {
	// Add validation to login
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// Connect db
		admin := &interfaces.AdminOnly{}
		if database.DB.Where("username = ?", username).First(&admin).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}

		// verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}
		return prepareResponse(admin, nil, true)
	} 
		return map[string]interface{}{"message": "not valid values"}
}

func GetAllUser(id string, jwt string) map[string]interface{} {

	if !helpers.ValidateAdminToken(jwt) {
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

func DeleteUser(id string, jwt string) map[string]interface{} {
	if !helpers.ValidateAdminToken(jwt) {
		return map[string]interface{}{"message": "Invalid token"}
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return map[string]interface{}{"message": "Invalid account ID"}
	}

	// Delete Account
	result := database.DB.Where("id = ?", idUint).Delete(&interfaces.User{})
	if result.RowsAffected == 0 {
		return map[string]interface{}{"message": "User not found or already deleted"}
	}
	return map[string]interface{}{"message": "User deleted successfully"}
}

func DeleteAccout(id string, jwt string) map[string]interface{} {
	if !helpers.ValidateAdminToken(jwt) {
		return map[string]interface{}{"message": "Invalid token"}
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return map[string]interface{}{"message": "Invalid account ID"}
	}

	// Delete Account
	result := database.DB.Where("id = ?", idUint).Delete(&interfaces.Account{})
	if result.RowsAffected == 0 {
		return map[string]interface{}{"message": "Account not found or already deleted"}
	}
	return map[string]interface{}{"message": "Account deleted successfully"}
}