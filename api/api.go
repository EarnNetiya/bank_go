package api

import (
	"encoding/json"
	"fmt"
	"goproject-bank/admin"
	"goproject-bank/helpers"
	"goproject-bank/transactions"
	"goproject-bank/useraccounts"
	"goproject-bank/users"
	"io/ioutil"
	"log"
	"net/http"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type Register struct {
	Username string
	Email    string
	Password string
}

type TransactionBody struct {
	UserId uint
	From   uint
	To     uint
	Amount int
}

type ErrResponse struct {
	Message string
}

func readBody(r *http.Request) ([]byte) {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "all is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		resp := call
		json.NewEncoder(w).Encode(resp)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// Read body
	body := readBody(r)
	// Handle Login
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)
	// Prepare response
	apiResponse(login, w)
}


func register(w http.ResponseWriter, r *http.Request) {
	// Read body
	body := readBody(r)
	// Handle registration
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	// Prepare response
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	user := users.GetUser(userId, auth)
	apiResponse(user, w)
}

func getMyTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userID"]
	auth := r.Header.Get("Authorization")

	transactions := transactions.GetMyTransactions(userId, auth)
	apiResponse(transactions, w)
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	auth := r.Header.Get("Authorization")
	// Handle registration
	var formattedBody TransactionBody
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	transaction := useraccounts.Transactions(formattedBody.UserId, formattedBody.From, formattedBody.To, formattedBody.Amount, auth)
	// Prepare response
	apiResponse(transaction, w)
}

// admin
func adminLogin(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)

	login := admin.Login(formattedBody.Username, formattedBody.Password)
	apiResponse(login, w)
}

func adminRegister(w http.ResponseWriter, r *http.Request) {
	// Read body
	body := readBody(r)
	// Handle registration
	var formattedBody Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	register := admin.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	// Prepare response
	apiResponse(register, w)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	admin_Id := vars["id"]
	auth := r.Header.Get("Authorization")

	response := admin.GetAllUser(admin_Id, auth)
	apiResponse(response, w)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user_Id := vars["id"]
	auth := r.Header.Get("Authorization")

	response := admin.DeleteUser(user_Id, auth)
	apiResponse(response, w)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acc_Id := vars["id"]
	auth := r.Header.Get("Authorization")

	response := admin.DeleteAccount(acc_Id, auth)
	apiResponse(response, w)
}


func StartApi() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/transactions", Transactions).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/transaction/{userID}", getMyTransactions).Methods("GET")

	// AdminOnly
	router.HandleFunc("/admin/login", adminLogin).Methods("POST")
	router.HandleFunc("/admin/register", adminRegister).Methods("POST")
	router.HandleFunc("/admin/user/{id}", getAllUsers).Methods("GET")
	router.HandleFunc("/delete/user/{user_Id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/delete/account/{acc_Id}", deleteAccount).Methods("DELETE")
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
