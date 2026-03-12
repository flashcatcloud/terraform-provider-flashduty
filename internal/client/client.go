package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultBaseURL = "https://api.flashcat.cloud"
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	BaseURL    string
	AppKey     string
	HTTPClient *http.Client
	UserAgent  string
}

type ClientOption func(*Client)

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = url
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func NewClient(appKey string, version string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		AppKey:  appKey,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		UserAgent: fmt.Sprintf("terraform-provider-flashduty/%s", version),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// APIError represents an error returned by the Flashduty API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Flashduty API error: %s - %s", e.Code, e.Message)
}

// doRequest performs an HTTP request and returns the response body.
func (c *Client) doRequest(ctx context.Context, method, path string, body any) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	u, err := url.Parse(fmt.Sprintf("%s%s", c.BaseURL, path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	q := u.Query()
	q.Set("app_key", c.AppKey)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, httpStatusToError(resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// doRequestWithResponse performs an HTTP request and unmarshals the response into the provided type.
func doRequestWithResponse[T any](c *Client, ctx context.Context, method, path string, body any) (*T, string, error) {
	respBody, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		return nil, "", err
	}

	var result struct {
		RequestID string    `json:"request_id,omitempty"`
		Error     *APIError `json:"error,omitempty"`
		Data      *T        `json:"data,omitempty"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, "", fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(respBody))
	}

	if result.Error != nil {
		return nil, result.RequestID, result.Error
	}

	return result.Data, result.RequestID, nil
}
