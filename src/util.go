package main

import (
	"net/http"
)

// HTTPGet is a helper function that automatically adds the
// tool's UA to an HTTP GET request
func HTTPGet(url string) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "e6dl: go edition (@tjhorner on Telegram)")

	return client.Do(req)
}
