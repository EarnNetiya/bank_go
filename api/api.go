package api

import (
	"encoding/json"
	"fmt"
	"goproject-bank/helpers"
	"goproject-bank/users"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type ErrResponse struct {
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {
	// Ready body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleErr(err)

	// Handle login
	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleErr(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)

	// Prepare response
	if login["message"] == "all is fine" {
		resp := login
		json.NewEncoder(w).Encode(resp)
	} else {
		// hAndle success
		resp := ErrResponse{Message: "Wrong username or password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func StartApi() {
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
