package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/saurabhmittal16/fluttermate/creds"
	"github.com/saurabhmittal16/fluttermate/firebase"
	"github.com/saurabhmittal16/fluttermate/score"
)

var clientSecret string
var clientID string

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

func getURL(id string) string {
	return fmt.Sprintf("https://api.github.com/user/%s?client_id=%s&client_secret=%s", id, clientID, clientSecret)
}

func checkAuth(f authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		data, err := firebase.VerifyToken(token)
		checkHTTPError(w, err, "Could not verify token", http.StatusInternalServerError)

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
	err := jsonResponse(w, r, text{
		Message: "Welcome to FlutterMate",
	})
	checkHTTPError(w, err, "Could not generate JSON Response", http.StatusInternalServerError)
}

func signupResponse(w http.ResponseWriter, r *http.Request, user tokenData) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// if user already exists, respond with login successful
	if firebase.DocExists(user.uid, "users") {
		fmt.Println("User exists, login successful")
		err := jsonResponse(w, r, authMessage{
			Code:    1,
			Message: "Login successful",
		})
		checkHTTPError(w, err, "Could not generate JSON Response", http.StatusInternalServerError)

	} else {
		// get user info from Github API
		profile, err := http.Get(getURL(user.ghID))
		checkHTTPError(w, err, "Failed to fetch data from Github API", http.StatusInternalServerError)

		// decode user from API response
		var newUser firebaseUser
		err = json.NewDecoder(profile.Body).Decode(&newUser)
		checkHTTPError(w, err, "Error while decoding API response", http.StatusInternalServerError)
		newUser.Email = user.email

		// save user to Firestore
		err = firebase.CreateUser(user.uid, "users", newUser)
		checkHTTPError(w, err, "Error while decoding API response", http.StatusInternalServerError)
		fmt.Println("User created successfuly")

		// respond with signup successful
		err = jsonResponse(w, r, authMessage{
			Code:    2,
			Message: "Signup successful",
		})
		checkHTTPError(w, err, "Could not generate JSON Response", http.StatusInternalServerError)

		// create go routine that calculates score and updates to firestore
		go func(ghID string, docID string) {
			start := time.Now()
			score, err := score.GetScore(ghID)
			if err != nil {
				fmt.Printf("Could not calculate score: %v\n", err)
				return
			}

			err = firebase.UpdateScore(docID, "users", score)
			if err != nil {
				fmt.Printf("Could not calculate score: %v\n", err)
				return
			}
			fmt.Printf("Score of user %s is %v, updated successfuly in %v\n", docID, score, time.Now().Sub(start))
		}(user.ghID, user.uid)
	}
}

func init() {
	/*
		pass path of JSON file containing Github app credentials of form-
		{
			"client_id": "xxxxx",
			"client_secret": "yyyyy"
		}
	*/
	clientID, clientSecret = creds.GetClientCreds("./github.json")
}

func main() {
	// Register routes
	http.HandleFunc("/", welcomeResponse)
	http.HandleFunc("/signup", checkAuth(signupResponse))

	// Start the server and log errors
	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
