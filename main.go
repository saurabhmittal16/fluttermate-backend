package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// jsonResponse generates JSON Response from given interface
func jsonResponse(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}

type text struct {
	Message string `json:"message"`
}

// welcomeResponse handles requests to root
func welcomeResponse(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, r, text{
		Message: "Welcome to FlutterMate",
	})
}

func main() {
	http.HandleFunc("/", welcomeResponse)

	fmt.Println("Server running at port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
