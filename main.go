package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var app *firebase.App
var client *firestore.Client
var err error

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
	// Open connection to Firestore
	client, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatal("Failed to intialise Firestore", err)
	}
	defer client.Close()

	readUserData()

	// Register routes
	http.HandleFunc("/", welcomeResponse)

	// Start the server and log errors
	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
