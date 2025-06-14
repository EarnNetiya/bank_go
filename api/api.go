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
	Username string `json:"username"`
	Password string `json:"password"`
}

type Register struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	InitialAmount int    `json:"initialAmount"` // ชื่อ field ตรงกับ JSON body
}

type TransactionBody struct {
	UserId uint `json:"userId"`
	From   uint `json:"from"`
	To     uint `json:"to"`
	Amount int  `json:"amount"`
}

func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	return body, err
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(call)
}

func login(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var formattedBody Login
	if err := json.Unmarshal(body, &formattedBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := users.Login(formattedBody.Username, formattedBody.Password)
	apiResponse(resp, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var formattedBody Register
	if err := json.Unmarshal(body, &formattedBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := users.Register(
		formattedBody.Username,
		formattedBody.Email,
		formattedBody.Password,
		formattedBody.InitialAmount,
	)
	apiResponse(resp, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	auth := r.Header.Get("Authorization")

	resp := users.GetUser(userId, auth)
	apiResponse(resp, w)
}

func getMyTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userID"]
	auth := r.Header.Get("Authorization")

	resp := transactions.GetMyTransactions(userId, auth)
	apiResponse(resp, w)
}

func Transactions(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	auth := r.Header.Get("Authorization")

	var formattedBody TransactionBody
	if err := json.Unmarshal(body, &formattedBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := useraccounts.Transactions(formattedBody.UserId, formattedBody.From, formattedBody.To, formattedBody.Amount, auth)
	apiResponse(resp, w)
}

// admin
func adminLogin(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	login := admin.Login(formattedBody.Username, formattedBody.Password)
	apiResponse(login, w)
}

func adminRegister(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var formattedBody Register
	err = json.Unmarshal(body, &formattedBody)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	register := admin.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
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
	user_Id := vars["user_Id"]
	auth := r.Header.Get("Authorization")

	response := admin.DeleteUser(user_Id, auth)
	apiResponse(response, w)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acc_Id := vars["acc_Id"]
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
