package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// initializeAppWithServiceAccount initialises the firebase admin sdk
func initializeAppWithServiceAccount() *firebase.App {
	opt := option.WithCredentialsFile("./service-account.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}

// jsonResponse generates JSON Response from given interface
func jsonResponse(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

// checkHTTPError logs and returns the error
func checkHTTPError(w http.ResponseWriter, err error, message string, code int) {
	if err != nil {
		fmt.Println(message, err)
		http.Error(w, message, code)
		return
	}
}
