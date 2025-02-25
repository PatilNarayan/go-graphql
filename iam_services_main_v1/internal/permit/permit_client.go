package permit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ResponseType represents the type of response expected from an API call
type ResponseType int

const (
	// MapResponse indicates a map[string]interface{} response type
	MapResponse ResponseType = iota
	// RawResponse indicates a generic interface{} response type
	RawResponse
)

// PermitClient is a client for interacting with the Permit API.
type PermitClient struct {
	baseURL string
	headers map[string]string
	client  *http.Client
}

// Config holds configuration for the PermitClient
type Config struct {
	PDPEndpoint string
	ProjectID   string
	EnvID       string
	APIKey      string
	Timeout     time.Duration
}

// NewPermitClient initializes a new PermitClient with default configuration from environment variables.
func NewPermitClient() *PermitClient {
	config := Config{
		PDPEndpoint: os.Getenv("PERMIT_PDP_ENDPOINT"),
		ProjectID:   os.Getenv("PERMIT_PROJECT"),
		EnvID:       os.Getenv("PERMIT_ENV"),
		APIKey:      os.Getenv("PERMIT_TOKEN"),
		Timeout:     30 * time.Second,
	}
	return NewPermitClientWithConfig(config)
}

// NewPermitClientWithConfig initializes a new PermitClient with the provided configuration.
func NewPermitClientWithConfig(config Config) *PermitClient {
	baseURL := fmt.Sprintf("%s/v2/facts/%s/%s", config.PDPEndpoint, config.ProjectID, config.EnvID)

	return &PermitClient{
		baseURL: baseURL,
		headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", config.APIKey),
			"Content-Type":  "application/json",
		},
		client: &http.Client{Timeout: config.Timeout},
	}
}

// SendRequest sends an HTTP request and returns the response as a map[string]interface{}.
func (pc *PermitClient) SendRequest(ctx context.Context, method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	result, err := pc.sendRequestWithType(ctx, method, endpoint, payload, MapResponse)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result.(map[string]interface{}), nil
}

// SendRequestInterface sends an HTTP request and returns the response as an interface{}.
func (pc *PermitClient) SendRequestInterface(ctx context.Context, method, endpoint string, payload interface{}) (interface{}, error) {
	return pc.sendRequestWithType(ctx, method, endpoint, payload, RawResponse)
}

// sendRequestWithType is the underlying implementation for both request methods.
func (pc *PermitClient) sendRequestWithType(ctx context.Context, method, endpoint string, payload interface{}, responseType ResponseType) (interface{}, error) {
	// Serialize payload to JSON
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal payload: %v", err)
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Determine the correct base URL
	url := pc.getURLForEndpoint(endpoint)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", url, endpoint), body)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Log URL
	log.Printf("permit request URL: %s", req.URL.String())

	// Add headers
	for key, value := range pc.headers {
		req.Header.Set(key, value)
	}

	// Send the request
	resp, err := pc.client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	return pc.processResponse(resp, responseType)
}

// getURLForEndpoint determines the correct base URL for the given endpoint
func (pc *PermitClient) getURLForEndpoint(endpoint string) string {
	// Use schema URL for roles/resources endpoints
	if strings.Contains(endpoint, "roles") || strings.Contains(endpoint, "resources") {
		return strings.Replace(pc.baseURL, "facts", "schema", 1)
	}
	return pc.baseURL
}

// processResponse handles the HTTP response and returns the appropriate result
func (pc *PermitClient) processResponse(resp *http.Response, responseType ResponseType) (interface{}, error) {
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("HTTP error: %d - %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Parse response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if len(respBody) == 0 {
		log.Printf("Empty response body")
		return nil, nil
	}

	var result interface{}
	if responseType == MapResponse {
		result = make(map[string]interface{})
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}
