package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/iterator"
)

// read user data
func readUserData() {
	iter := fsClient.Collection("users").Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		for k, v := range doc.Data() {
			fmt.Printf("%s -> %v\n", k, v)
		}
		fmt.Println()
	}
}

// find if user exists
func userExists(uid string) bool {
	dsnap, err := fsClient.Collection("users").Doc(uid).Get(context.Background())
	if err != nil {
		fmt.Println("Checking of document failed", err)
		return false
	}

	return dsnap.Exists()
}

// create user
func createUser(uid string, v interface{}) error {
	_, err := fsClient.Collection("users").Doc(uid).Set(context.Background(), v)
	return err
}
