package score

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type language struct {
	Dart int `json:"Dart"`
}

type repo struct {
	Language string `json:"language"`
	IsFork   bool   `json:"fork"`
	Forks    int    `json:"forks"`
	Stars    int    `json:"stargazers_count"`
	Watchers int    `json:"watchers_count"`
	URL      string `json:"languages_url"`
}

type repos []repo

func must(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getURL(id string) string {
	return fmt.Sprintf("https://api.github.com/user/%s/repos", id)
}

// GetScore takes GitHub User ID as input
// and returns the Flutter score of the user
func GetScore(id string) (float64, error) {
	resp, err := http.Get(getURL(id))
	if err != nil {
		return -1, err
	}

	var result, filtered repos
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return -1, err
	}

	for _, res := range result {
		if res.Language == "Dart" && res.IsFork == false {
			filtered = append(filtered, res)
		}
	}

	lineChan := make(chan float64)
	score := 0.0

	for _, r := range filtered {
		score += float64(r.Forks + r.Stars + r.Watchers)

		go func(url string) {
			var temp language
			resp, err := http.Get(url)
			must(err)

			err = json.NewDecoder(resp.Body).Decode(&temp)
			must(err)

			lineChan <- float64(temp.Dart) / 100
		}(r.URL)
	}

	for range filtered {
		score += <-lineChan
	}

	return score, nil
}
