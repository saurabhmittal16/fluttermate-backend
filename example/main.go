package main

import (
	"fmt"
	"log"
	"time"

	"github.com/saurabhmittal16/fluttermate/creds"
	"github.com/saurabhmittal16/fluttermate/score"
)

func main() {
	clientID, clientSecret := creds.GetClientCreds("../github.json")
	score.Init(clientID, clientSecret)

	start := time.Now()

	fscore, err := score.GetScore("31154435")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got score %v in time %v\n", fscore, time.Now().Sub(start))
}
