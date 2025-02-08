package permit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// PermitClient is a client for interacting with the Permit API.
type PermitClient struct {
	BaseURL string
	Headers map[string]string
	Client  *http.Client
}

// NewPermitClient initializes a new PermitClient.
func NewPermitClient() *PermitClient {
	baseURL := os.Getenv("PERMIT_PDP_ENDPOINT")
	projectID := os.Getenv("PERMIT_PROJECT")
	envID := os.Getenv("PERMIT_ENV")
	apiKey := os.Getenv("PERMIT_TOKEN")
	return &PermitClient{
		BaseURL: fmt.Sprintf("%s/v2/facts/%s/%s", baseURL, projectID, envID),
		Headers: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", apiKey),
			"Content-Type":  "application/json",
		},
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

// sendRequest sends an HTTP request and handles retries.
func (pc *PermitClient) sendRequest(ctx context.Context, method, endpoint string, payload interface{}) (interface{}, error) {
	var result interface{}

	operation := func() error {
		// Serialize payload to JSON
		var body io.Reader
		if payload != nil {
			jsonData, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Failed to marshal payload: %v", err)
				return backoff.Permanent(err)
			}
			body = bytes.NewBuffer(jsonData)
		}

		if strings.Contains(endpoint, "roles") || strings.Contains(endpoint, "resources") {
			pc.BaseURL = strings.Replace(pc.BaseURL, "facts", "schema", 1)
		}
		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", pc.BaseURL, endpoint), body)
		if err != nil {
			log.Printf("Failed to create HTTP request: %v", err)
			return backoff.Permanent(err)
		}

		//add log url
		log.Printf("permit request URL: %s", req.URL.String())

		// Add headers
		for key, value := range pc.Headers {
			req.Header.Set(key, value)
		}

		// Send the request
		resp, err := pc.Client.Do(req)
		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return err
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("HTTP error: %d - %s", resp.StatusCode, string(body))
			return errors.New(fmt.Sprintf("HTTP error: %d", resp.StatusCode))
		}

		// Parse response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return backoff.Permanent(err)
		}

		if len(respBody) == 0 {
			log.Printf("Empty response body")
			return nil
		}

		if err := json.Unmarshal(respBody, &result); err != nil {
			log.Printf("Failed to unmarshal response: %v", err)
			return backoff.Permanent(err)
		}

		return nil
	}

	// Use exponential backoff for retries
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 30 * time.Second

	if err := backoff.Retry(operation, bo); err != nil {
		log.Printf("Request failed after retries: %v", err)
		return nil, err
	}

	return result, nil
}

// // CreateResource creates a resource with the given data.
// func (pc *PermitClient) CreateResource(ctx context.Context, resourceData map[string]interface{}) (map[string]interface{}, error) {
// 	return pc.sendRequest(ctx, "POST", "resources", resourceData)
// }

// // CreateTenant creates a tenant with the given data.
// func (pc *PermitClient) CreateTenant(ctx context.Context, tenantData map[string]interface{}) (map[string]interface{}, error) {
// 	return pc.sendRequest(ctx, "POST", "tenants", tenantData)
// }

// // DeleteTenant deletes a tenant by ID.
// func (pc *PermitClient) DeleteTenant(ctx context.Context, tenantID string) error {
// 	_, err := pc.sendRequest(ctx, "DELETE", fmt.Sprintf("tenants/%s", tenantID), nil)
// 	return err
// }

// // CreateResourceInstance creates a resource instance with the given data.
// func (pc *PermitClient) CreateResourceInstance(ctx context.Context, instanceData map[string]interface{}) (map[string]interface{}, error) {
// 	return pc.sendRequest(ctx, "POST", "resource_instances", instanceData)
// }

// func (pc *PermitClient) UpdateTenant(ctx context.Context, tenantID string, tenantName string) (map[string]interface{}, error) {
// 	return pc.sendRequest(ctx, "PATCH", fmt.Sprintf("tenants/%s", tenantID), map[string]interface{}{"name": tenantName})
// }

func (pc *PermitClient) SendRequest(ctx context.Context, method, endpoint string, payload interface{}) (interface{}, error) {
	return pc.sendRequest(ctx, method, endpoint, payload)
}
