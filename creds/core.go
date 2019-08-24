package creds

import (
	"encoding/json"
	"log"
	"os"
)

type cred struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

// GetClientCreds returns the client ID and client secret for the Github Application
func GetClientCreds(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("No credentials file available: %v\n", err)
	}

	var temp cred
	err = json.NewDecoder(file).Decode(&temp)

	if err != nil {
		log.Fatalf("Could not decode JSON: %v\n", err)
	}

	return temp.ID, temp.Secret
}
