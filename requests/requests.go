// requests.go
package requests

import (
	"appjet-client-cli/configurations"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var tokenFilePath = "token.txt"

type LoginResponse struct {
	AccessToken string `json:"ACCESS_TOKEN"`
}

func setAccessToken(token string) {
	err := ioutil.WriteFile(tokenFilePath, []byte(token), 0644)
	if err != nil {
		fmt.Println("Error saving token to file:", err)
	}
}

func getAccessToken() string {
	token, err := ioutil.ReadFile(tokenFilePath)
	if err != nil {
		return ""
	}
	return string(token)
}

func addAuthorizationHeader(req *http.Request) {
	token := getAccessToken()
	if token != "" {
		req.Header.Set("Authorization", token)
	}
}

func DoLogin(username, password string) {
	loginURL := configurations.AppConfig.AppJetURL + "/api/login"

	requestBody := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var loginResponse LoginResponse
		err := json.Unmarshal(body, &loginResponse)
		if err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return
		}

		token := loginResponse.AccessToken
		setAccessToken(token)
		fmt.Println("Login successful. Access Token:", token)
	} else {
		fmt.Printf("Login failed. Server response: %s\n", string(body))
	}
}

func DeleteLoginState() int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", configurations.AppConfig.AppJetURL+"/api/logout", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}

	addAuthorizationHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Logout successful, delete the token file
		err := deleteTokenFile()
		if err != nil {
			fmt.Println("Error deleting token file:", err)
		}

		fmt.Println("Logout done. Thank you for using APPJET.")
		return 0
	} else {
		fmt.Println("Logout failed. Status code:", resp.StatusCode)
		return 1
	}
}

func deleteTokenFile() error {
	err := ioutil.WriteFile(tokenFilePath, []byte(""), 0644)
	return err
}

func SignupUser(username, password, email string) {
	payload := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)

	url := configurations.AppConfig.AppJetURL + "/api/signup"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	addAuthorizationHeader(req)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Account created successfully!")
	} else {
		fmt.Println("Failed to create account. Status code:", resp.StatusCode)
	}
}
