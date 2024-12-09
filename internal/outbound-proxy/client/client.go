package client

import "net/http"

var client = &http.Client{}

func Fetch(req *http.Request) (*http.Response, error) {

	var resp *http.Response

	resp, err := client.Do(req)
	if err != nil {
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
	}

	return nil, nil
}
