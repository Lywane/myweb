package http

import (
	"io"
	"time"
	"io/ioutil"
	"bytes"
	"net/http"
)

const DefaultTimeout = time.Duration(10 * time.Second)

func get(url string, timeout ...time.Duration) ([]byte, error) {
	return request(http.MethodGet, url, nil, timeout...)
}

func post(url string, body []byte, timeout ...time.Duration) ([]byte, error) {
	data := bytes.NewReader(body)
	return request(http.MethodPost, url, data, timeout...)
}

func request(method, url string, body io.Reader, timeout ...time.Duration) ([]byte, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte(""), err
	}
	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")
	}
	duration := DefaultTimeout
	if len(timeout) > 1 {
		duration = timeout[0]
	}
	client := http.Client{Timeout: duration}
	resp, err := client.Do(request)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return ret, nil
}
