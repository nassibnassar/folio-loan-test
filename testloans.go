package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	limit := os.Args[1]

	err := run(limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}

func run(limit string) error {
	token, err := login()
	if err != nil {
		return err
	}

	err = retrieveLoans(token, limit)
	if err != nil {
		return err
	}

	return nil
}

func login() (string, error) {
	client := &http.Client{}

	m := map[string]string{
		"username": "admin",
		"password": "admin",
	}
	json, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST",
		"http://localhost:9130/authn/login",
		bytes.NewBuffer(json))
	if err != nil {
		return "", err
	}

	req.Header.Add("X-Okapi-Tenant", "diku")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json,text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Header["X-Okapi-Token"][0], nil
}

func retrieveLoans(token string, limit string) error {
	client := &http.Client{}

	url := "http://localhost:9130/loan-storage/loans?limit=" +
		limit +
		"&offset=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-Okapi-Tenant", "diku")
	req.Header.Add("X-Okapi-Token", token)
	req.Header.Add("Accept", "application/json,text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", body)
	fmt.Println("URL:", url)
	fmt.Println("Status code:", resp.StatusCode,
		http.StatusText(resp.StatusCode))

	return nil
}
