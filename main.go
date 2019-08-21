package main

import (
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
)

var app *firebase.App

type text struct {
	Message string `json:"message"`
}

// welcomeResponse handles requests to root
func welcomeResponse(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, text{
		Message: "Welcome to FlutterMate",
	})
}

func init() {
	app = initializeAppWithServiceAccount()
}

func main() {
	http.HandleFunc("/", welcomeResponse)

	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
