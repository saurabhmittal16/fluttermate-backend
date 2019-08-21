package main

import (
	"context"
	"encoding/json"
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

type authHandler func(http.ResponseWriter, *http.Request, tokenData)

type text struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type authMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type tokenData struct {
	uid   string
	ghID  string
	email string
}

// ToDo: Set email from token and store GitHub ID as string in Firestore
type firebaseUser struct {
	Name      string `json:"name" firestore:"name"`
	Picture   string `json:"avatar_url" firestore:"picture"`
	Email     string `json:"email" firestore:"email"`
	Username  string `json:"login" firestore:"username"`
	Location  string `json:"location" firestore:"location"`
	Bio       string `json:"bio" firestore:"bio"`
	Repos     int    `json:"public_repos" firestore:"repos"`
	Followers int    `json:"followers" firestore:"followers"`
	Following int    `json:"following" firestore:"following"`
	Github    int    `json:"id" firestore:"github"`
}

func checkAuth(f authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		data, err := authClient.VerifyIDToken(context.Background(), token)
		if err != nil {
			fmt.Printf("Could not verify token: %v\n", err)
			http.Error(w, "Could not verify token", http.StatusInternalServerError)
			return
		}
		temp := data.Claims["firebase"].(map[string]interface{})["identities"].(map[string]interface{})
		user := tokenData{
			uid:   data.UID,
			ghID:  temp["github.com"].([]interface{})[0].(string),
			email: temp["email"].([]interface{})[0].(string),
		}
		// fmt.Println(user)
		f(w, r, user)
	}
}

// welcomeResponse handles requests to root
func welcomeResponse(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, text{
		Message: "Welcome to FlutterMate",
	})
}

func profileResponse(w http.ResponseWriter, r *http.Request, user tokenData) {
	jsonResponse(w, r, text{
		Message: fmt.Sprintf("Hello, %s!", user.email),
	})
}

// ToDo: Extract error handling to seperate function
func signupResponse(w http.ResponseWriter, r *http.Request, user tokenData) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// if user already exists, respond with login successful
	if userExists(user.uid) {
		fmt.Println("User exists, login successful")
		jsonResponse(w, r, authMessage{
			Code:    1,
			Message: "Login successful",
		})
	} else {
		// get user info from Github API
		profile, err := http.Get("https://api.github.com/user/" + user.ghID)
		if err != nil {
			http.Error(w, "Failed to fetch data from Github API", http.StatusInternalServerError)
			return
		}

		// decode user from API response
		var newUser firebaseUser
		err = json.NewDecoder(profile.Body).Decode(&newUser)
		if err != nil {
			fmt.Println("Error while decoding API response", err)
			http.Error(w, "Error while decoding API response", http.StatusInternalServerError)
			return
		}

		// save user to Firestore
		fmt.Printf("User is %+v\n", newUser)
		err = createUser(user.uid, newUser)
		if err != nil {
			fmt.Println("Error while decoding API response", err)
			http.Error(w, "Error while decoding API response", http.StatusInternalServerError)
			return
		}

		// respond with signup successful
		jsonResponse(w, r, authMessage{
			Code:    2,
			Message: "Signup successful",
		})
	}
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
	http.HandleFunc("/signup", checkAuth(signupResponse))

	// Start the server and log errors
	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
