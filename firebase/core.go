package firebase

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var app *firebase.App
var fsClient *firestore.Client
var authClient *auth.Client
var err error

func init() {
	// Initialise firebase app
	opt := option.WithCredentialsFile("./service-account.json")
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Initialise firestore client
	fsClient, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to intialise Firestore: %v\n", err)
	}

	// Initialise auth client
	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("Error getting Auth client: %v\n", err)
	}
}

// VerifyToken decodes the token
func VerifyToken(token string) (*auth.Token, error) {
	return authClient.VerifyIDToken(context.Background(), token)
}

// ReadUserData reads data from a collection
func ReadUserData(collection string) {
	iter := fsClient.Collection(collection).Documents(context.Background())
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

// DocExists finds whether given uid is in the collection
func DocExists(uid string, collection string) bool {
	dsnap, err := fsClient.Collection(collection).Doc(uid).Get(context.Background())
	if err != nil {
		fmt.Println("Checking of document failed", err)
		return false
	}

	return dsnap.Exists()
}

// CreateUser creates a new collection
func CreateUser(uid string, collection string, v interface{}) error {
	_, err := fsClient.Collection(collection).Doc(uid).Set(context.Background(), v)
	return err
}

// UpdateScore accepts user ID and score and updates DB
func UpdateScore(uid string, collection string, score float64) error {
	_, err := fsClient.Collection(collection).Doc(uid).Set(context.Background(), map[string]interface{}{
		"score": score,
	}, firestore.MergeAll)

	return err
}
