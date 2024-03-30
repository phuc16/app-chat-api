package curl

import (
	"io"
	"net/http"
)

func Get(url string, body io.Reader) (*http.Response, error) {
	return curl("GET", url, body)
}

func Post(url string, body io.Reader) (*http.Response, error) {
	return curl("POST", url, body)
}

func curl(method, url string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	return client.Do(req)
}
