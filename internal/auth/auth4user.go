package auth

import (
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/term"
)

// GetUserInput prompts the user for input and returns it
func GetUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// GetPasswordInput prompts the user for password input and returns it
func GetPasswordInput(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		return ""
	}
	fmt.Println()
	return strings.TrimSpace(string(bytePassword))
}

// AuthenticateUser authenticates the user with the given username and password
func AuthenticateUser(username, password string) (string, string, error) {
	if database.GetUsersCollection() == nil {
		return "", "", fmt.Errorf("users collection is not initialized")
	}

	token, role, err := loginUser(username, password)
	if err != nil {
		return "", "", err
	}

	return token, role, nil
}

// loginUser sends a login request and returns the token and role
func loginUser(username, password string) (string, string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := serializer.JSONSerializerInstance.Serialize(loginData)
	if err != nil {
		return "", "", err
	}

	resp, err := http.Post("http://localhost:8080/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("login request failed with status code: %d", resp.StatusCode)
		}
		return "", "", fmt.Errorf(string(body))
	}

	return parseLoginResponse(resp.Body)
}

// parseLoginResponse parses the login response and returns the token and role
func parseLoginResponse(body io.Reader) (string, string, error) {
	var result map[string]string
	err := serializer.JSONSerializerInstance.DeserializeFromReader(body, &result)
	if err != nil {
		return "", "", err
	}

	token, ok := result["token"]
	if !ok {
		return "", "", fmt.Errorf("failed to get token from login response")
	}

	role, ok := result["role"]
	if !ok {
		return "", "", fmt.Errorf("failed to get role from login response")
	}

	return token, role, nil
}
