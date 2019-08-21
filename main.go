package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

var app *firebase.App
var fsClient *firestore.Client
var authClient *auth.Client
var err error

type text struct {
	Message string `json:"message"`
}

type tokenData struct {
	ghID  string
	email string
}

func checkAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		data, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			// ToDo: Handle this gracefully
			log.Fatalf("Could not verify token: %v\n", err)
		}
		temp := data.Claims["firebase"].(map[string]interface{})["identities"].(map[string]interface{})
		user := tokenData{
			ghID:  temp["github.com"].([]interface{})[0].(string),
			email: temp["email"].([]interface{})[0].(string),
		}
		fmt.Println(user)
		f(w, r)
	}
}

// welcomeResponse handles requests to root
func welcomeResponse(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, text{
		Message: "Welcome to FlutterMate",
	})
}

func profileResponse(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, text{
		Message: "Hello, World!",
	})
}

func init() {
	app = initializeAppWithServiceAccount()
}

func main() {
	// Open connection to Firestore
	fsClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to intialise Firestore: %v\n", err)
	}
	defer fsClient.Close()

	// Create auth client
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Auth client: %v\n", err)
	}

	// Register routes
	http.HandleFunc("/", welcomeResponse)
	http.HandleFunc("/me", checkAuth(profileResponse))

	// Start the server and log errors
	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
