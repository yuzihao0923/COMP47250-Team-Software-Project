package client

import (
	"net/http"
)

// GetClientWithToken returns an HTTP client with the provided JWT token set in the Authorization header
func GetClientWithToken(token string) *http.Client {
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
