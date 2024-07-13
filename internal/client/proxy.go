package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Get available broker
func GetBroker(proxyURL string) (*Broker, error) {
	resp, err := http.Get(fmt.Sprintf("%s/get-broker", proxyURL))
	if err != nil {
		return nil, fmt.Errorf("failed to get broker from proxy: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get broker, status code: %d, message: %s", resp.StatusCode, string(body))
	}

	var broker Broker
	if err := json.NewDecoder(resp.Body).Decode(&broker); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &broker, nil
}
