package api

import (
	"encoding/json"
	"fmt"
	"goproject-bank/admin"
	"goproject-bank/blockchain"
	"goproject-bank/database"
	"goproject-bank/interfaces"
	"strconv"

	// "goproject-bank/database"
	"goproject-bank/helpers"
	// "goproject-bank/interfaces"
	"goproject-bank/transactions"
	"goproject-bank/users"
	"io/ioutil"
	"log"
	"net/http"

	// "strconv"

	// "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	// "github.com/jinzhu/gorm"
)

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type Register struct {
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Balance int    `json:"Balance"` 
}

type TransactionBody struct {
    FromAccountNumber string `json:"fromAccountNumber"`
    ToAccountNumber   string `json:"toAccountNumber"`
    Amount            int    `json:"amount"`
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
		log.Println("Failed to read request body:", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var formattedBody UserLogin
	if err := json.Unmarshal(body, &formattedBody); err != nil {
		log.Println("Invalid JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	resp := users.Login(formattedBody.Email, formattedBody.Password)
	apiResponse(resp, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Println("Register request body:", string(body))

	var formattedBody Register
	if err := json.Unmarshal(body, &formattedBody); err != nil {
		log.Println("Invalid JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Register request parsed: %+v\n", formattedBody)

	resp := users.Register(
		formattedBody.Email,
		formattedBody.Username,
		formattedBody.Password,
		formattedBody.Balance,
	)
	log.Println("Register response:", resp)

	apiResponse(resp, w)
}


func getMyTransactions(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["userID"]
    token := helpers.ExtractTokenFromRequest(r)

    log.Println("Token:", token)

    result := transactions.GetMyTransactions(userID, token)
    json.NewEncoder(w).Encode(result)
}

func getMyAccounts(w http.ResponseWriter, r *http.Request) {
    tokenString := helpers.ExtractTokenFromRequest(r)
    userID, err := helpers.ExtractUserID(tokenString)
    if err != nil {
        log.Println("Invalid token:", err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var accounts []interfaces.Account
    if err := database.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
        log.Println("Error fetching accounts:", err)
        http.Error(w, "Error fetching accounts", http.StatusInternalServerError)
        return
    }

    var respAccounts []interfaces.ResponseAccount
    for _, acc := range accounts {
        respAccounts = append(respAccounts, interfaces.ResponseAccount{
            ID:            acc.ID,
            Name:          acc.Name,
            Balance:       int(acc.Balance),
            AccountNumber: acc.AccountNumber,
        })
    }

    response := map[string]interface{}{
        "message": "All accounts retrieved",
        "data":    respAccounts,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func Transactions(w http.ResponseWriter, r *http.Request) {
    tokenString := helpers.ExtractTokenFromRequest(r)
    if tokenString == "" {
        log.Println("Missing token in request")
        http.Error(w, "Missing token", http.StatusUnauthorized)
        return
    }

    userID, err := helpers.ExtractUserID(tokenString)
    if err != nil {
        log.Println("Token validation failed:", err)
        http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
        return
    }

    var req TransactionBody
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        log.Println("Invalid request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    log.Println("Received transaction request:", req) // Debug

    result := transactions.CreateTransactionByAccountNumbers(req.FromAccountNumber, req.ToAccountNumber, req.Amount, userID)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// admin
func adminLogin(w http.ResponseWriter, r *http.Request) {
	body, err := readBody(r)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	var formattedBody AdminLogin
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

func getAccount(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accID := vars["id"]
    auth := r.Header.Get("Authorization")

    log.Println("Authorization Header:", auth)

    resp := admin.GetAccount(accID, auth)
    if status, ok := resp["status"].(int); ok {
        w.WriteHeader(status)
    } else {
        w.WriteHeader(http.StatusOK)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}


func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["id"]
    auth := r.Header.Get("Authorization")

    log.Println("Authorization Header:", auth)

    resp := admin.GetUser(userID, auth)
    if status, ok := resp["status"].(int); ok {
        w.WriteHeader(status)
    } else {
        w.WriteHeader(http.StatusOK)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	response := admin.GetAllUsers(auth)
	apiResponse(response, w)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_Id"]
    auth := r.Header.Get("Authorization")
    log.Println("Received Authorization header:", auth)
    response := admin.DeleteUser(userID, auth)
    if status, ok := response["status"].(int); ok {
        w.WriteHeader(status)
    }
    apiResponse(response, w)
}

func deleteAccount(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    accID := vars["acc_Id"]
    auth := r.Header.Get("Authorization")
    response := admin.DeleteAccount(accID, auth)
    if status, ok := response["status"].(int); ok {
        w.WriteHeader(status)
    } else {
        w.WriteHeader(http.StatusOK)
	}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// blockchain Admin
func getBlockchainTransactions(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["id"]

    userIDUint, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        log.Printf("Invalid user ID: %v", err)
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    auth := r.Header.Get("Authorization")
    if !helpers.ValidateAdminToken(auth) {
        log.Printf("Unauthorized access attempt for user %s", userIDStr)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

	var accounts []interfaces.Account
	if err := database.DB.Where("user_id = ?", userIDUint).Find(&accounts).Error; err != nil {
		log.Printf("Error fetching accounts for user %d: %v", userIDUint, err)
		http.Error(w, "Error fetching accounts", http.StatusInternalServerError)
		return
	}
	
	log.Printf("Found %d accounts for user %d", len(accounts), userIDUint)
	if len(accounts) == 0 {
		log.Printf("No accounts found for user %s, returning empty data.", userIDStr) // Added log
		response := map[string]interface{}{
			"message": "No accounts found for user " + userIDStr,
			"data":    []interfaces.BlockWithHash{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	
	accountSet := make(map[string]bool)
	for _, acc := range accounts {
		accountSet[acc.AccountNumber] = true
		log.Printf("Account number for user %d: %s", userIDUint, acc.AccountNumber)
	}
	
	chainWithHashes := blockchain.Chain.GetBlockchainWithHashes()
	log.Printf("Total blocks in blockchain: %d", len(chainWithHashes))
	
	var userTransactions []interfaces.BlockWithHash
	for i, block := range chainWithHashes { // Added index for better logging
		log.Printf("Checking block %d: Sender=%s, Receiver=%s", i, block.Data.SenderAccount, block.Data.ReceiverAccount) // Added log
		if accountSet[block.Data.SenderAccount] || accountSet[block.Data.ReceiverAccount] {
			userTransactions = append(userTransactions, block)
			log.Printf("Transaction found for user %d - Block Index: %d, Sender: %s, Receiver: %s", userIDUint, i, block.Data.SenderAccount, block.Data.ReceiverAccount)
		}
	}

    log.Printf("Total transactions found for user %d: %d", userIDUint, len(userTransactions))

    response := map[string]interface{}{
        "message": "Blockchain transactions for user " + userIDStr,
        "data":    userTransactions,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func StartApi() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/transactions", Transactions).Methods("POST")
	router.HandleFunc("/transaction/{userID}", getMyTransactions).Methods("GET")
	router.HandleFunc("/myaccount", getMyAccounts).Methods("GET")

	// AdminOnly
	router.HandleFunc("/admin/login", adminLogin).Methods("POST")
	router.HandleFunc("/admin/register", adminRegister).Methods("POST")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/account/{id}", getAccount).Methods("GET")
	router.HandleFunc("/admin/users", getAllUsers).Methods("GET")
	router.HandleFunc("/delete/user/{user_Id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/delete/account/{acc_Id}", deleteAccount).Methods("DELETE")
	router.HandleFunc("/admin/blockchain/{id}", getBlockchainTransactions).Methods("GET")
	
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
