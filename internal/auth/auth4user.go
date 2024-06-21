package auth

import (
	"COMP47250-Team-Software-Project/internal/database"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/term"
)

func GetUserInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

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

func AuthenticateUser(username, password string) (string, string, error) {
	if database.GetUsersCollection() == nil {
		return "", "", fmt.Errorf("users collection is not initialized")
	}

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
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

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
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
