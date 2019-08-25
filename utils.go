package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/saurabhmittal16/fluttermate/firebase"
	"github.com/saurabhmittal16/fluttermate/score"
)

// jsonResponse generates JSON Response from given interface
func jsonResponse(w http.ResponseWriter, r *http.Request, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(v)
	return err
}

// checkHTTPError logs and returns the error
func checkHTTPError(w http.ResponseWriter, err error, message string, code int) {
	if err != nil {
		fmt.Println(message, err)
		http.Error(w, message, code)
		return
	}
}

func seedData(users []tokenData) {
	for _, user := range users {
		// get user info from Github API
		profile, err := http.Get(getURL(user.ghID))
		if err == nil {
			fmt.Println("Received data from API")
		}

		// decode user from API response
		var newUser firebaseUser
		err = json.NewDecoder(profile.Body).Decode(&newUser)
		newUser.Email = user.email
		if err == nil {
			fmt.Println("Data decoded successfuly")
		}

		// save user to Firestore
		err = firebase.CreateUser(user.uid, "users", newUser)
		fmt.Println("User created successfuly", err)

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
	// ToDo: handle program termination before go routine finishes
}
