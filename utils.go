package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
