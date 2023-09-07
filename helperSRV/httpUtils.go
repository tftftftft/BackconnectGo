package main

import (
	"bytes"
	"io"
	"net/http"
)

func ReadRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	return body, err
}

func ForwardRequest(body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", "http://38.180.61.234:20001/addProxy", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}
