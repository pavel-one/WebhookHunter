package main

import (
	"io"
	"net/http"
	"testing"
)

func TestMultiCreateHunter(t *testing.T) {
	for i := 0; i < 1000; i++ {
		requestCreateHunter(t)
	}
}

func requestCreateHunter(t *testing.T) {
	url := "http://localhost:3000/api/v1/"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		t.Fatal(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(body))
}
