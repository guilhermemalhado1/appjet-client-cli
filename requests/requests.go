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

type LoginResponse struct {
	AccessToken string `json:"ACCESS_TOKEN"`
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
		fmt.Println("Login successful. Access Token:", token)
	} else {
		fmt.Printf("Login failed. Server response: %s\n", string(body))
	}
}


func DeleteLoginState() int {
	resp, err := http.Get(configurations.AppConfig.AppJetURL + "/api/logout")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return 0
	} else {
		fmt.Println("Logout failed. Status code:", resp.StatusCode)
	}

	return 1
}

func SignupUser(username, password, email string) {
	payload := fmt.Sprintf(`{"username": "%s", "password": "%s", "email": "%s"}`, username, password, email)

	url := configurations.AppConfig.AppJetURL + "/api/signup"

	resp, err := http.Post(url, "application/json", strings.NewReader(payload))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Account created successfully! Now you must click on the activation link in your email to be able to login.")
	} else {
		fmt.Println("Failed to create account. Status code:", resp.StatusCode)
	}
}
