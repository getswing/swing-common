package sw

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func DoRequest(ctx context.Context, method, url string, body interface{}, headers map[string]string, maxRetries int) ([]byte, int, error) {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Default headers
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response

	// Retry loop
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, 0, fmt.Errorf("request failed after %d attempts: %w", attempt, err)
			}
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond) // exponential backoff
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
		}

		if resp.StatusCode >= 500 && attempt < maxRetries {
			// Retry on 5xx
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
			continue
		}

		return respBody, resp.StatusCode, nil
	}

	return nil, 0, errors.New("unexpected error in DoRequest")
}
