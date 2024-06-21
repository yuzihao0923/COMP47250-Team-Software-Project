package api

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
)

func getJWTToken() (string, error) {
    loginURL := "http://localhost:8080/login" // 确保这个URL是正确的
    creds := map[string]string{
        "username": "user",    // 你的用户名
        "password": "password", // 你的密码
    }
    jsonData, err := json.Marshal(creds)
    if err != nil {
        return "", fmt.Errorf("error marshalling credentials: %v", err)
    }

    resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("error sending login request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to login, status code: %d", resp.StatusCode)
    }

    var result map[string]string
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        return "", fmt.Errorf("error decoding response: %v", err)
    }

    token, ok := result["token"]
    if !ok {
        return "", fmt.Errorf("token not found in response")
    }

    return token, nil
}

func getClientWithToken() *http.Client {
    token := os.Getenv("JWT_TOKEN")
    if token == "" {
        var err error
        token, err = getJWTToken()
        if err != nil {
            panic(fmt.Sprintf("Failed to get JWT token: %v", err))
        }
        os.Setenv("JWT_TOKEN", token) // 设置环境变量，便于后续使用
    }

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
