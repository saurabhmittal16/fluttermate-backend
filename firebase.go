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
