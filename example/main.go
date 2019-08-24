package main

import (
	"fmt"
	"log"
	"time"

	"github.com/saurabhmittal16/fluttermate/score"
)

func main() {
	start := time.Now()

	fscore, err := score.GetScore("34417814")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fscore)
	fmt.Println(time.Now().Sub(start))
}
