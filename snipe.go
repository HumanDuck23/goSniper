package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getToken(email string, password string) string {
	type ResJSON struct {
		ClientToken string `json:"clientToken"`
		AccessToken string `json:"accessToken"`
	}

	url := "https://authserver.mojang.com/authenticate"
	var jsonStr = []byte(`{"username":"` + email + `", "password": "` + password + `"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode == 200 {
		// Logged in!
		success("Logged in!")
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			error("Error reading access token!")
			os.Exit(-1)
		}
		var rj ResJSON
		json.Unmarshal([]byte(string(body)), &rj)
		return rj.AccessToken
	} else {
		fmt.Println(res)
	}

	return ""
}

func validateToken(token string) bool {
	url := "https://authserver.mojang.com/validate"
	var jsonStr = []byte(`{"accessToken":"` + token + `"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	return res.StatusCode == 204
}