package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var httpClient = http.Client{Timeout: 5 * time.Second}

// GetJSON will send get request to ai backend, and parse json response
func GetJSON(url string, target interface{}) error {
	resp, err := httpClient.Get(withEndpoint(url))
	if err != nil {
		return err
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("get failed: %w", err)
	}

	err = json.Unmarshal(b, target)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	return nil
}

// PostJSON will send post request to ai backend, and parse json response
func PostJSON(url string, body interface{}, target interface{}) error {
	var bodyReader io.Reader = nil
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body failed: %w", err)
		}
		bodyReader = bytes.NewBuffer(b)
	}

	resp, err := httpClient.Post(withEndpoint(url), "application/json", bodyReader)
	if err != nil {
		return fmt.Errorf("post failed: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, target)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	return nil
}

// withEndpoint returns full url with endpoint
func withEndpoint(url string) string {
	if strings.HasPrefix(url, "/") {
		return fmt.Sprintf("%s%s", endpoint, url)
	}
	return url
}
