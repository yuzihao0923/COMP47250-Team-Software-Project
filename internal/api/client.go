package api

import (
	"bytes"
	"fmt"
	"net/http"

	"COMP47250-Team-Software-Project/pkg/serializer"
)

func GetJWTToken(username, password string) (string, error) {
	loginURL := "http://localhost:8080/login"
	creds := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := serializer.JSONSerializerInstance.Serialize(creds)
	if err != nil {
		fmt.Println("Error marshalling credentials:", err)
		return "", fmt.Errorf("error marshalling credentials: %v", err)
	}

	fmt.Println("Sending login request to", loginURL)
	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending login request:", err)
		return "", fmt.Errorf("error sending login request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Login request sent, status code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
	}

	var result map[string]string
	err = serializer.JSONSerializerInstance.DeserializeFromReader(resp.Body, &result)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	token, ok := result["token"]
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	fmt.Println("Received JWT token:", token)
	return token, nil
}

func getClientWithToken(token string) *http.Client {
	client := &http.Client{}
	client.Transport = &transportWithToken{token, http.DefaultTransport}
	return client
}

type transportWithToken struct {
	token     string
	transport http.RoundTripper
}

func (t *transportWithToken) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.transport.RoundTrip(req)
}
