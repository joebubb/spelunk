package util

import "net/http"

func UrlIsValidGet(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}

	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		return true
	default:
		return false
	}
}
