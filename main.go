package easynet

import (
	"io/ioutil"
	"log"
	"net/http"
)

func Easyget(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error performing GET request: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	return string(body), nil
}
